package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SystemMetricsModel = (*customSystemMetricsModel)(nil)

type (
	// SystemMetricsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSystemMetricsModel.
	SystemMetricsModel interface {
		systemMetricsModel
		withSession(session sqlx.Session) SystemMetricsModel
	}

	customSystemMetricsModel struct {
		*defaultSystemMetricsModel
	}
)

// NewSystemMetricsModel returns a model for the database table.
func NewSystemMetricsModel(conn sqlx.SqlConn) SystemMetricsModel {
	return &customSystemMetricsModel{
		defaultSystemMetricsModel: newSystemMetricsModel(conn),
	}
}

func (m *customSystemMetricsModel) withSession(session sqlx.Session) SystemMetricsModel {
	return NewSystemMetricsModel(sqlx.NewSqlConnFromSession(session))
}
