package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ InstancePortMappingModel = (*customInstancePortMappingModel)(nil)

type (
	// InstancePortMappingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInstancePortMappingModel.
	InstancePortMappingModel interface {
		instancePortMappingModel
		withSession(session sqlx.Session) InstancePortMappingModel
	}

	customInstancePortMappingModel struct {
		*defaultInstancePortMappingModel
	}
)

// NewInstancePortMappingModel returns a model for the database table.
func NewInstancePortMappingModel(conn sqlx.SqlConn) InstancePortMappingModel {
	return &customInstancePortMappingModel{
		defaultInstancePortMappingModel: newInstancePortMappingModel(conn),
	}
}

func (m *customInstancePortMappingModel) withSession(session sqlx.Session) InstancePortMappingModel {
	return NewInstancePortMappingModel(sqlx.NewSqlConnFromSession(session))
}
