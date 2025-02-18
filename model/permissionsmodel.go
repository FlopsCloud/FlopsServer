package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PermissionsModel = (*customPermissionsModel)(nil)

type (
	// PermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPermissionsModel.
	PermissionsModel interface {
		permissionsModel
		WithSession(session sqlx.Session) PermissionsModel
		FindByFilter(ctx context.Context, name, code string, page, pageSize uint64) ([]*Permissions, uint64, error)
		FindOneByCode(ctx context.Context, code string) (*Permissions, error)
	}

	customPermissionsModel struct {
		*defaultPermissionsModel
	}
)

// NewPermissionsModel returns a model for the database table.
func NewPermissionsModel(conn sqlx.SqlConn) PermissionsModel {
	return &customPermissionsModel{
		defaultPermissionsModel: newPermissionsModel(conn),
	}
}

func (m *customPermissionsModel) WithSession(session sqlx.Session) PermissionsModel {
	return NewPermissionsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customPermissionsModel) FindByFilter(ctx context.Context, name, code string, page, pageSize uint64) ([]*Permissions, uint64, error) {
	where := "1=1"
	var args []interface{}
	if name != "" {
		where += " AND permission_name LIKE ?"
		args = append(args, "%"+name+"%")
	}
	if code != "" {
		where += " AND permission_name = ?"
		args = append(args, code)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, where)
	var total uint64
	err := m.conn.QueryRowCtx(ctx, &total, query, args...)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*Permissions{}, 0, nil
	}

	query = fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY permission_id DESC LIMIT ?,?", permissionsRows, m.table, where)
	args = append(args, (page-1)*pageSize, pageSize)
	var resp []*Permissions
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

func (m *customPermissionsModel) FindOneByCode(ctx context.Context, code string) (*Permissions, error) {
	var resp Permissions
	query := fmt.Sprintf("select %s from %s where permission_name = ? limit 1", permissionsRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, code)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
