package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ServerDiscountsModel = (*customServerDiscountsModel)(nil)

type (
	// ServerDiscountsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customServerDiscountsModel.
	ServerDiscountsModel interface {
		serverDiscountsModel
		withSession(session sqlx.Session) ServerDiscountsModel
		Count(ctx context.Context, orgId, serverId uint64) (uint64, error)
		FindMany(ctx context.Context, orgId, serverId uint64, offset, limit uint64) ([]*ServerDiscounts, error)
	}

	customServerDiscountsModel struct {
		*defaultServerDiscountsModel
	}
)

// NewServerDiscountsModel returns a model for the database table.
func NewServerDiscountsModel(conn sqlx.SqlConn) ServerDiscountsModel {
	return &customServerDiscountsModel{
		defaultServerDiscountsModel: newServerDiscountsModel(conn),
	}
}

func (m *customServerDiscountsModel) withSession(session sqlx.Session) ServerDiscountsModel {
	return NewServerDiscountsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customServerDiscountsModel) Count(ctx context.Context, orgId, serverId uint64) (uint64, error) {
	var count uint64
	query := fmt.Sprintf("select count(*) from %s where `is_deleted` = 0", m.table)

	if orgId > 0 {
		query += fmt.Sprintf(" and `org_id` = %d", orgId)
	}
	if serverId > 0 {
		query += fmt.Sprintf(" and `server_id` = %d", serverId)
	}

	err := m.conn.QueryRowCtx(ctx, &count, query)
	return count, err
}

func (m *customServerDiscountsModel) FindMany(ctx context.Context, orgId, serverId uint64, offset, limit uint64) ([]*ServerDiscounts, error) {
	query := fmt.Sprintf("select %s from %s where `is_deleted` = 0", serverDiscountsRows, m.table)

	if orgId > 0 {
		query += fmt.Sprintf(" and `org_id` = %d", orgId)
	}
	if serverId > 0 {
		query += fmt.Sprintf(" and `server_id` = %d", serverId)
	}

	query += " order by `srv_discount_id` desc limit ?,?"

	var resp []*ServerDiscounts
	err := m.conn.QueryRowsCtx(ctx, &resp, query, offset, limit)
	return resp, err
}
