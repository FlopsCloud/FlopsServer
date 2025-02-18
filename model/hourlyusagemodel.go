package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ HourlyUsageModel = (*customHourlyUsageModel)(nil)

type (
	// HourlyUsageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHourlyUsageModel.
	HourlyUsageModel interface {
		hourlyUsageModel
		WithSession(session sqlx.Session) HourlyUsageModel
		FindLastUncharged(ctx context.Context) (*HourlyUsage, error)
		FindLatestByInstanceId(ctx context.Context, orgId uint64, userId uint64, instanceId uint64) (*HourlyUsage, error)
		FindUnchargedComplete(ctx context.Context, orgId uint64, userId uint64, instanceId uint64) ([]*HourlyUsage, error)
		List(ctx context.Context, usageId, orgId, userId *uint64, startDatetime, endDatetime *uint64, page, pageSize uint64) ([]*HourlyUsageView, uint64, error)
	}

	customHourlyUsageModel struct {
		*defaultHourlyUsageModel
	}
)

// NewHourlyUsageModel returns a model for the database table.
func NewHourlyUsageModel(conn sqlx.SqlConn) HourlyUsageModel {
	return &customHourlyUsageModel{
		defaultHourlyUsageModel: newHourlyUsageModel(conn),
	}
}

func (m *customHourlyUsageModel) WithSession(session sqlx.Session) HourlyUsageModel {
	return NewHourlyUsageModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customHourlyUsageModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, fn)
}

func (m *defaultHourlyUsageModel) FindLatestByInstanceId(ctx context.Context, orgId uint64, userId uint64, instanceId uint64) (*HourlyUsage, error) {
	query := fmt.Sprintf("select * from %s where org_id = ? and user_id = ? and instance_id = ?  order by hournum desc limit 1", m.table)
	var resp HourlyUsage
	err := m.conn.QueryRowCtx(ctx, &resp, query, orgId, userId, instanceId)
	return &resp, err
}

func (m *defaultHourlyUsageModel) FindUnchargedComplete(ctx context.Context, orgId uint64, userId uint64, instanceId uint64) ([]*HourlyUsage, error) {
	query := fmt.Sprintf("select * from %s where org_id = ? and user_id = ? and instance_id = ? and is_charged = 0 and minute_end = 60", m.table)
	var resp []*HourlyUsage
	err := m.conn.QueryRowsCtx(ctx, &resp, query, orgId, userId, instanceId)
	return resp, err
}

func (m *defaultHourlyUsageModel) FindLastUncharged(ctx context.Context) (*HourlyUsage, error) {
	query := fmt.Sprintf("select * from %s where is_charged = 0 order by usage_datetime desc limit 1", m.table)
	var resp HourlyUsage
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

type HourlyUsageView struct {
	HourlyUsage
	InstanceName string `db:"instance_name"`
	ResourceName string `db:"resource_name"`
}

func (m *customHourlyUsageModel) List(ctx context.Context, usageId, orgId, userId *uint64, startDatetime, endDatetime *uint64, page, pageSize uint64) ([]*HourlyUsageView, uint64, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}

	if usageId != nil {
		conditions = append(conditions, "a.`usage_id` = ?")
		args = append(args, *usageId)
	}
	if orgId != nil {
		conditions = append(conditions, "a.`org_id` = ?")
		args = append(args, *orgId)
	}
	if userId != nil {
		conditions = append(conditions, "a.`user_id` = ?")
		args = append(args, *userId)
	}
	if startDatetime != nil {
		conditions = append(conditions, "a.`usage_datetime` >= ?")
		args = append(args, *startDatetime)
	}
	if endDatetime != nil {
		conditions = append(conditions, "a.`usage_datetime` <= ?")
		args = append(args, *endDatetime)
	}

	whereClause := strings.Join(conditions, " AND ")

	whereClauseCount := strings.Replace(whereClause, "a.`", "`", -1)

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, whereClauseCount)
	var total uint64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * pageSize
	fields := []string{}
	for _, field := range hourlyUsageFieldNames {
		fields = append(fields, fmt.Sprintf("a.%s", field))
	}
	allfields := strings.Join(fields, ",")

	query := fmt.Sprintf(`SELECT %s, ifnull(instances.name,"") as instance_name, ifnull(running_resources.resource_name,"") as resource_name 
	 FROM %s a 
LEFT JOIN 
    instances ON a.instance_id = instances.instance_id 
LEFT JOIN 
    running_resources ON a.runres_id = running_resources.runres_id  
WHERE %s ORDER BY a.usage_date DESC LIMIT ? OFFSET ?`, allfields, m.table, whereClause)
	args = append(args, pageSize, offset)

	var resp []*HourlyUsageView
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
