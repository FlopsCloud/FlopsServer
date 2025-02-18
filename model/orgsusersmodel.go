package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OrgsUsersModel = (*customOrgsUsersModel)(nil)

type (
	// OrgsUsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrgsUsersModel.
	OrgsUsersModel interface {
		orgsUsersModel
		WithSession(session sqlx.Session) OrgsUsersModel
	}

	customOrgsUsersModel struct {
		*defaultOrgsUsersModel
	}
)

// NewOrgsUsersModel returns a model for the database table.
func NewOrgsUsersModel(conn sqlx.SqlConn) OrgsUsersModel {
	return &customOrgsUsersModel{
		defaultOrgsUsersModel: newOrgsUsersModel(conn),
	}
}

func (m *customOrgsUsersModel) WithSession(session sqlx.Session) OrgsUsersModel {
	return NewOrgsUsersModel(sqlx.NewSqlConnFromSession(session))
}
