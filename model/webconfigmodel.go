package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WebConfigModel = (*customWebConfigModel)(nil)

type (
	// WebConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWebConfigModel.
	WebConfigModel interface {
		webConfigModel
		withSession(session sqlx.Session) WebConfigModel
	}

	customWebConfigModel struct {
		*defaultWebConfigModel
	}
)

// NewWebConfigModel returns a model for the database table.
func NewWebConfigModel(conn sqlx.SqlConn) WebConfigModel {
	return &customWebConfigModel{
		defaultWebConfigModel: newWebConfigModel(conn),
	}
}

func (m *customWebConfigModel) withSession(session sqlx.Session) WebConfigModel {
	return NewWebConfigModel(sqlx.NewSqlConnFromSession(session))
}
