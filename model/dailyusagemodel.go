package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DailyUsageModel = (*customDailyUsageModel)(nil)

type (
	// DailyUsageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDailyUsageModel.
	DailyUsageModel interface {
		dailyUsageModel
		withSession(session sqlx.Session) DailyUsageModel
		List(ctx context.Context, orgId, userId *uint64, startDate, endDate *string, page, pageSize uint64) ([]*DailyUsage, uint64, error)
	}

	customDailyUsageModel struct {
		*defaultDailyUsageModel
	}
)

// NewDailyUsageModel returns a model for the database table.
func NewDailyUsageModel(conn sqlx.SqlConn) DailyUsageModel {
	return &customDailyUsageModel{
		defaultDailyUsageModel: newDailyUsageModel(conn),
	}
}

func (m *customDailyUsageModel) withSession(session sqlx.Session) DailyUsageModel {
	return NewDailyUsageModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customDailyUsageModel) List(ctx context.Context, orgId, userId *uint64, startDate, endDate *string, page, pageSize uint64) ([]*DailyUsage, uint64, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}

	if orgId != nil {
		conditions = append(conditions, "`org_id` = ?")
		args = append(args, *orgId)
	}
	if userId != nil {
		conditions = append(conditions, "`user_id` = ?")
		args = append(args, *userId)
	}
	if startDate != nil {
		conditions = append(conditions, "`usage_date` >= ?")
		args = append(args, *startDate)
	}
	if endDate != nil {
		conditions = append(conditions, "`usage_date` <= ?")
		args = append(args, *endDate)
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
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY usage_date DESC LIMIT ? OFFSET ?", dailyUsageRows, m.table, whereClause)
	args = append(args, pageSize, offset)

	var resp []*DailyUsage
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
