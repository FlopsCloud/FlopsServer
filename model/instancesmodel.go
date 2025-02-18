package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ InstancesModel = (*customInstancesModel)(nil)

type (
	// InstancesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInstancesModel.
	InstancesModel interface {
		instancesModel
		WithSession(session sqlx.Session) InstancesModel
		FindByFilter(ctx context.Context, userId, serverId uint64, page, pageSize uint64) ([]*Instances, uint64, error)
	}

	customInstancesModel struct {
		*defaultInstancesModel
	}
)

// NewInstancesModel returns a model for the database table.
func NewInstancesModel(conn sqlx.SqlConn) InstancesModel {
	return &customInstancesModel{
		defaultInstancesModel: newInstancesModel(conn),
	}
}

func (m *customInstancesModel) WithSession(session sqlx.Session) InstancesModel {
	return NewInstancesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customInstancesModel) FindByFilter(ctx context.Context, userId, serverId uint64, page, pageSize uint64) ([]*Instances, uint64, error) {
	where := "1=1"
	var args []interface{}

	if userId > 0 {
		where += " AND user_id = ?"
		args = append(args, userId)
	}

	if serverId > 0 {
		where += " AND server_id = ?"
		args = append(args, serverId)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, where)
	var total uint64
	err := m.conn.QueryRowCtx(ctx, &total, query, args...)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*Instances{}, 0, nil
	}

	query = fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY instance_id DESC LIMIT ?,?", instancesRows, m.table, where)
	args = append(args, (page-1)*pageSize, pageSize)
	var resp []*Instances
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
