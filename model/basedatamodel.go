package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BaseDataModel = (*customBaseDataModel)(nil)

type (
	// BaseDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBaseDataModel.
	BaseDataModel interface {
		baseDataModel
		withSession(session sqlx.Session) BaseDataModel
	}

	customBaseDataModel struct {
		*defaultBaseDataModel
	}
)

// NewBaseDataModel returns a model for the database table.
func NewBaseDataModel(conn sqlx.SqlConn) BaseDataModel {
	return &customBaseDataModel{
		defaultBaseDataModel: newBaseDataModel(conn),
	}
}

func (m *customBaseDataModel) withSession(session sqlx.Session) BaseDataModel {
	return NewBaseDataModel(sqlx.NewSqlConnFromSession(session))
}
