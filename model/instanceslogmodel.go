package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ InstancesLogModel = (*customInstancesLogModel)(nil)

type (
	// InstancesLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInstancesLogModel.
	InstancesLogModel interface {
		instancesLogModel
		withSession(session sqlx.Session) InstancesLogModel
	}

	customInstancesLogModel struct {
		*defaultInstancesLogModel
	}
)

// NewInstancesLogModel returns a model for the database table.
func NewInstancesLogModel(conn sqlx.SqlConn) InstancesLogModel {
	return &customInstancesLogModel{
		defaultInstancesLogModel: newInstancesLogModel(conn),
	}
}

func (m *customInstancesLogModel) withSession(session sqlx.Session) InstancesLogModel {
	return NewInstancesLogModel(sqlx.NewSqlConnFromSession(session))
}
