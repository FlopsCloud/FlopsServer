package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ServerOrgsModel = (*customServerOrgsModel)(nil)

type (
	// ServerOrgsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customServerOrgsModel.
	ServerOrgsModel interface {
		serverOrgsModel
		withSession(session sqlx.Session) ServerOrgsModel
	}

	customServerOrgsModel struct {
		*defaultServerOrgsModel
	}
)

// NewServerOrgsModel returns a model for the database table.
func NewServerOrgsModel(conn sqlx.SqlConn) ServerOrgsModel {
	return &customServerOrgsModel{
		defaultServerOrgsModel: newServerOrgsModel(conn),
	}
}

func (m *customServerOrgsModel) withSession(session sqlx.Session) ServerOrgsModel {
	return NewServerOrgsModel(sqlx.NewSqlConnFromSession(session))
}
