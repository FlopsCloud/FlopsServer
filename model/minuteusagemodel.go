package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MinuteUsageModel = (*customMinuteUsageModel)(nil)

type (
	// MinuteUsageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMinuteUsageModel.
	MinuteUsageModel interface {
		minuteUsageModel
		withSession(session sqlx.Session) MinuteUsageModel
		List(ctx context.Context, usageId, orgId, userId *uint64, startDatetime, endDatetime *string, page, pageSize uint64) ([]*MinuteUsage, uint64, error)
	}

	customMinuteUsageModel struct {
		*defaultMinuteUsageModel
	}
)

// NewMinuteUsageModel returns a model for the database table.
func NewMinuteUsageModel(conn sqlx.SqlConn) MinuteUsageModel {
	return &customMinuteUsageModel{
		defaultMinuteUsageModel: newMinuteUsageModel(conn),
	}
}

func (m *customMinuteUsageModel) withSession(session sqlx.Session) MinuteUsageModel {
	return NewMinuteUsageModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customMinuteUsageModel) List(ctx context.Context, usageId, orgId, userId *uint64, startDatetime, endDatetime *string, page, pageSize uint64) ([]*MinuteUsage, uint64, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}

	if usageId != nil {
		conditions = append(conditions, "`usage_id` = ?")
		args = append(args, *usageId)
	}
	if orgId != nil {
		conditions = append(conditions, "`org_id` = ?")
		args = append(args, *orgId)
	}
	if userId != nil {
		conditions = append(conditions, "`user_id` = ?")
		args = append(args, *userId)
	}
	if startDatetime != nil {
		conditions = append(conditions, "`usage_datetime` >= ?")
		args = append(args, *startDatetime)
	}
	if endDatetime != nil {
		conditions = append(conditions, "`usage_datetime` <= ?")
		args = append(args, *endDatetime)
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
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY usage_datetime DESC LIMIT ? OFFSET ?", minuteUsageRows, m.table, whereClause)
	args = append(args, pageSize, offset)

	var resp []*MinuteUsage
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
