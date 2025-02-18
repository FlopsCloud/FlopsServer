package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ResourceDiscountsModel = (*customResourceDiscountsModel)(nil)

type (
	// ResourceDiscountsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customResourceDiscountsModel.
	ResourceDiscountsModel interface {
		resourceDiscountsModel
		withSession(session sqlx.Session) ResourceDiscountsModel
		Count(ctx context.Context, orgId uint64) (int64, error)
		FindByConditions(ctx context.Context, orgId uint64, page, pageSize uint64) ([]*ResourceDiscounts, error)
	}

	customResourceDiscountsModel struct {
		*defaultResourceDiscountsModel
	}
)

// NewResourceDiscountsModel returns a model for the database table.
func NewResourceDiscountsModel(conn sqlx.SqlConn) ResourceDiscountsModel {
	return &customResourceDiscountsModel{
		defaultResourceDiscountsModel: newResourceDiscountsModel(conn),
	}
}

func (m *customResourceDiscountsModel) withSession(session sqlx.Session) ResourceDiscountsModel {
	return NewResourceDiscountsModel(sqlx.NewSqlConnFromSession(session))
}

// Count returns the total number of resource discounts for an organization
func (m *customResourceDiscountsModel) Count(ctx context.Context, orgId uint64) (int64, error) {
	var count int64
	condition := " 1 = 1 "

	if orgId > 0 {
		condition += fmt.Sprintf(" and org_id = %d ", orgId)
	}

	query := fmt.Sprintf("select count(*) from %s where %s", m.table, condition)
	err := m.conn.QueryRowCtx(ctx, &count, query)
	return count, err
}

// FindByConditions returns a list of resource discounts based on conditions with pagination
func (m *customResourceDiscountsModel) FindByConditions(ctx context.Context, orgId uint64, page, pageSize uint64) ([]*ResourceDiscounts, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	condition := " 1 = 1 "

	if orgId > 0 {
		condition += fmt.Sprintf(" and org_id = %d ", orgId)
	}

	query := fmt.Sprintf("select %s from %s where %s order by discount_id desc limit ?, ?", resourceDiscountsRows, m.table, condition)
	var resp []*ResourceDiscounts
	err := m.conn.QueryRowsCtx(ctx, &resp, query, offset, pageSize)
	return resp, err
}
