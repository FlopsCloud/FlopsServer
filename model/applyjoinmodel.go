package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ApplyJoinModel = (*customApplyJoinModel)(nil)

type (
	// ApplyJoinModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApplyJoinModel.
	ApplyJoinModel interface {
		applyJoinModel
		withSession(session sqlx.Session) ApplyJoinModel
	}

	customApplyJoinModel struct {
		*defaultApplyJoinModel
	}
)

// NewApplyJoinModel returns a model for the database table.
func NewApplyJoinModel(conn sqlx.SqlConn) ApplyJoinModel {
	return &customApplyJoinModel{
		defaultApplyJoinModel: newApplyJoinModel(conn),
	}
}

func (m *customApplyJoinModel) withSession(session sqlx.Session) ApplyJoinModel {
	return NewApplyJoinModel(sqlx.NewSqlConnFromSession(session))
}
