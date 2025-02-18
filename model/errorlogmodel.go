package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ErrorLogModel = (*customErrorLogModel)(nil)

type (
	// ErrorLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customErrorLogModel.
	ErrorLogModel interface {
		errorLogModel
		withSession(session sqlx.Session) ErrorLogModel
	}

	customErrorLogModel struct {
		*defaultErrorLogModel
	}
)

// NewErrorLogModel returns a model for the database table.
func NewErrorLogModel(conn sqlx.SqlConn) ErrorLogModel {
	return &customErrorLogModel{
		defaultErrorLogModel: newErrorLogModel(conn),
	}
}

func (m *customErrorLogModel) withSession(session sqlx.Session) ErrorLogModel {
	return NewErrorLogModel(sqlx.NewSqlConnFromSession(session))
}
