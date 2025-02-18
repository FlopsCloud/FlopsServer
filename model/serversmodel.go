package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ServersModel = (*customServersModel)(nil)

type (
	// ServersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customServersModel.
	ServersModel interface {
		serversModel
		withSession(session sqlx.Session) ServersModel
	}

	customServersModel struct {
		*defaultServersModel
	}
)

// NewServersModel returns a model for the database table.
func NewServersModel(conn sqlx.SqlConn) ServersModel {
	return &customServersModel{
		defaultServersModel: newServersModel(conn),
	}
}

func (m *customServersModel) withSession(session sqlx.Session) ServersModel {
	return NewServersModel(sqlx.NewSqlConnFromSession(session))
}
