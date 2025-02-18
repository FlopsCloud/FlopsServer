package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ServerTagsModel = (*customServerTagsModel)(nil)

type (
	// ServerTagsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customServerTagsModel.
	ServerTagsModel interface {
		serverTagsModel
		withSession(session sqlx.Session) ServerTagsModel
	}

	customServerTagsModel struct {
		*defaultServerTagsModel
	}
)

// NewServerTagsModel returns a model for the database table.
func NewServerTagsModel(conn sqlx.SqlConn) ServerTagsModel {
	return &customServerTagsModel{
		defaultServerTagsModel: newServerTagsModel(conn),
	}
}

func (m *customServerTagsModel) withSession(session sqlx.Session) ServerTagsModel {
	return NewServerTagsModel(sqlx.NewSqlConnFromSession(session))
}
