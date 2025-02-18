package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ResourcesModel = (*customResourcesModel)(nil)

type (
	// ResourcesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customResourcesModel.
	ResourcesModel interface {
		resourcesModel
		WithSession(session sqlx.Session) ResourcesModel
		List(ctx context.Context, orgId *uint64, resourceType *string, page, pageSize uint64) ([]*Resources, uint64, error)
	}

	customResourcesModel struct {
		*defaultResourcesModel
	}
)

// NewResourcesModel returns a model for the database table.
func NewResourcesModel(conn sqlx.SqlConn) ResourcesModel {
	return &customResourcesModel{
		defaultResourcesModel: newResourcesModel(conn),
	}
}

func (m *customResourcesModel) WithSession(session sqlx.Session) ResourcesModel {
	return NewResourcesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customResourcesModel) List(ctx context.Context, orgId *uint64, resourceType *string, page, pageSize uint64) ([]*Resources, uint64, error) {
	conditions := []string{"`is_deleted` = 0"}
	args := []interface{}{}

	if orgId != nil {
		conditions = append(conditions, "`org_id` = ?")
		args = append(args, *orgId)
	}
	if resourceType != nil {
		conditions = append(conditions, "`resource_type` = ?")
		args = append(args, *resourceType)
	}

	whereClause := strings.Join(conditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, whereClause)
	var total uint64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY created_at DESC LIMIT ? OFFSET ?", resourcesRows, m.table, whereClause)
	args = append(args, pageSize, offset)

	var resp []*Resources
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
