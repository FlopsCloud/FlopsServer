package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SystemNotificationsModel = (*customSystemNotificationsModel)(nil)

type (
	// SystemNotificationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSystemNotificationsModel.
	SystemNotificationsModel interface {
		systemNotificationsModel
		withSession(session sqlx.Session) SystemNotificationsModel
		FindByFilter(ctx context.Context, userId uint64, orgId uint64, typ string, page, pageSize int64) ([]*SystemNotifications, uint64, error)
		FindByFilter2(ctx context.Context, userId uint64, orgId uint64, typ string, status string, page, pageSize int64) ([]*SystemNotificationsEx, uint64, error)
	}

	customSystemNotificationsModel struct {
		*defaultSystemNotificationsModel
	}
)

// NewSystemNotificationsModel returns a model for the database table.
func NewSystemNotificationsModel(conn sqlx.SqlConn) SystemNotificationsModel {
	return &customSystemNotificationsModel{
		defaultSystemNotificationsModel: newSystemNotificationsModel(conn),
	}
}

func (m *customSystemNotificationsModel) withSession(session sqlx.Session) SystemNotificationsModel {
	return NewSystemNotificationsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSystemNotificationsModel) FindByFilter(ctx context.Context, userId uint64, orgId uint64, typ string, page, pageSize int64) ([]*SystemNotifications, uint64, error) {
	whereBuilder := strings.Builder{}
	args := []interface{}{}

	whereBuilder.WriteString("WHERE 1=1 ")

	if userId != 0 {
		whereBuilder.WriteString(" AND (user_id = ? OR user_id IS NULL)")
		args = append(args, userId)
	}

	if orgId != 0 {
		whereBuilder.WriteString(" AND (org_id = ? OR org_id IS NULL)")
		args = append(args, orgId)
	}

	if typ != "" {
		whereBuilder.WriteString(" AND `type` = ?")
		args = append(args, typ)
	}

	// Get total count
	var total uint64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereBuilder.String())
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get notifications with pagination
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY created_at DESC LIMIT ? OFFSET ?",
		m.table, whereBuilder.String())
	args = append(args, pageSize, offset)

	var notifications []*SystemNotifications
	err = m.conn.QueryRowsCtx(ctx, &notifications, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

type SystemNotificationsEx struct {
	NotificationId uint64        `db:"notification_id"`
	UserId         sql.NullInt64 `db:"user_id"` // 用户ID(NULL表示全局通知)
	OrgId          sql.NullInt64 `db:"org_id"`  // 组织ID(NULL表示不限组织)
	Title          string        `db:"title"`   // 通知标题
	Content        string        `db:"content"` // 通知内容
	Type           string        `db:"type"`    // 类型:system/maintenance/billing/security

	CreatedAt time.Time `db:"created_at"`
	IsRead    uint64    `db:"is_read"`
}

func (m *customSystemNotificationsModel) FindByFilter2(ctx context.Context, userId uint64, orgId uint64, typ string, status string, page, pageSize int64) ([]*SystemNotificationsEx, uint64, error) {

	// Get notifications with pagination
	offset := (page - 1) * pageSize

	//
	whereBuilder2 := strings.Builder{}
	args2 := []interface{}{}

	whereBuilder2.WriteString(" where 1=1")

	if userId != 0 {
		whereBuilder2.WriteString(" and (a.user_id = ? OR a.user_id is NULL)")
		args2 = append(args2, userId)
	}

	if orgId != 0 {
		whereBuilder2.WriteString("  and (a.org_id = ? OR a.org_id is NULL)")
		args2 = append(args2, orgId)
	}

	if typ != "" {
		whereBuilder2.WriteString(" AND a.`type` = ?")
		args2 = append(args2, typ)
	}

	if status == "unread" {
		whereBuilder2.WriteString(" AND IFNULL(b.is_read, 0) = 0")
	} else if status == "read" {
		whereBuilder2.WriteString(" AND IFNULL(b.is_read, 0) = 1")
	}

	// Get total count
	var total uint64
	countQuery := fmt.Sprintf("SELECT count(*)  FROM system_notifications a LEFT OUTER JOIN notification_reads b on a.notification_id=b.notification_id  %s ", whereBuilder2.String())
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args2...)
	if err != nil {
		return nil, 0, err
	}

	//IFNULL(column_a, 0) AS column_a_fixed
	query := fmt.Sprintf("SELECT a.notification_id,a.user_id,a.org_id,a.title,a.content,a.`type`,a.created_at,IFNULL(b.is_read, 0) AS is_read  FROM system_notifications a LEFT OUTER JOIN notification_reads b on a.notification_id=b.notification_id   %s ORDER BY a.created_at DESC LIMIT ? OFFSET ?", whereBuilder2.String())
	args2 = append(args2, pageSize, offset)

	var notifications []*SystemNotificationsEx
	err = m.conn.QueryRowsCtx(ctx, &notifications, query, args2...)
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}
