package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DiscountsModel = (*customDiscountsModel)(nil)

type (
	// DiscountsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDiscountsModel.
	DiscountsModel interface {
		discountsModel
		withSession(session sqlx.Session) DiscountsModel
		Count(ctx context.Context, orgId, resourceId uint64) (uint64, error)
		FindMany(ctx context.Context, orgId, resourceId uint64, offset, limit uint64) ([]*Discounts, error)
	}

	customDiscountsModel struct {
		*defaultDiscountsModel
	}
)

// NewDiscountsModel returns a model for the database table.
func NewDiscountsModel(conn sqlx.SqlConn) DiscountsModel {
	return &customDiscountsModel{
		defaultDiscountsModel: newDiscountsModel(conn),
	}
}

func (m *customDiscountsModel) withSession(session sqlx.Session) DiscountsModel {
	return NewDiscountsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customDiscountsModel) Count(ctx context.Context, orgId, resourceId uint64) (uint64, error) {
	var count uint64
	query := fmt.Sprintf("select count(*) from %s where `is_deleted` = 0", m.table)

	if orgId > 0 {
		query += fmt.Sprintf(" and `org_id` = %d", orgId)
	}
	if resourceId > 0 {
		query += fmt.Sprintf(" and `resource_id` = %d", resourceId)
	}

	err := m.conn.QueryRowCtx(ctx, &count, query)
	return count, err
}

func (m *customDiscountsModel) FindMany(ctx context.Context, orgId, resourceId uint64, offset, limit uint64) ([]*Discounts, error) {
	query := fmt.Sprintf("select %s from %s where `is_deleted` = 0", discountsRows, m.table)

	if orgId > 0 {
		query += fmt.Sprintf(" and `org_id` = %d", orgId)
	}
	if resourceId > 0 {
		query += fmt.Sprintf(" and `resource_id` = %d", resourceId)
	}

	query += " order by `discount_id` desc limit ?,?"

	var resp []*Discounts
	err := m.conn.QueryRowsCtx(ctx, &resp, query, offset, limit)
	return resp, err
}
