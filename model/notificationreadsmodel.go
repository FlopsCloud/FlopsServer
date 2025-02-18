package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ NotificationReadsModel = (*customNotificationReadsModel)(nil)

type (
	// NotificationReadsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotificationReadsModel.
	NotificationReadsModel interface {
		notificationReadsModel
		withSession(session sqlx.Session) NotificationReadsModel
	}

	customNotificationReadsModel struct {
		*defaultNotificationReadsModel
	}
)

// NewNotificationReadsModel returns a model for the database table.
func NewNotificationReadsModel(conn sqlx.SqlConn) NotificationReadsModel {
	return &customNotificationReadsModel{
		defaultNotificationReadsModel: newNotificationReadsModel(conn),
	}
}

func (m *customNotificationReadsModel) withSession(session sqlx.Session) NotificationReadsModel {
	return NewNotificationReadsModel(sqlx.NewSqlConnFromSession(session))
}
