package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ResourceOrgsModel = (*customResourceOrgsModel)(nil)

type (
	// ResourceOrgsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customResourceOrgsModel.
	ResourceOrgsModel interface {
		resourceOrgsModel
		withSession(session sqlx.Session) ResourceOrgsModel
		DeleteByUserIDOrgID(ctx context.Context, rid uint64, oid uint64) error
	}

	customResourceOrgsModel struct {
		*defaultResourceOrgsModel
	}
)

// NewResourceOrgsModel returns a model for the database table.
func NewResourceOrgsModel(conn sqlx.SqlConn) ResourceOrgsModel {
	return &customResourceOrgsModel{
		defaultResourceOrgsModel: newResourceOrgsModel(conn),
	}
}

func (m *customResourceOrgsModel) withSession(session sqlx.Session) ResourceOrgsModel {
	return NewResourceOrgsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customResourceOrgsModel) DeleteByUserIDOrgID(ctx context.Context, rid uint64, oid uint64) error {
	query := fmt.Sprintf("delete from %s where `resource_id` = ? and  `org_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, rid, oid)
	return err
}
