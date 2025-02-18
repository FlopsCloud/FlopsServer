package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ResourcesUserModel = (*customResourcesUserModel)(nil)

type (
	// ResourcesUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customResourcesUserModel.
	ResourcesUserModel interface {
		resourcesUserModel
		withSession(session sqlx.Session) ResourcesUserModel
	}

	customResourcesUserModel struct {
		*defaultResourcesUserModel
	}
)

// NewResourcesUserModel returns a model for the database table.
func NewResourcesUserModel(conn sqlx.SqlConn) ResourcesUserModel {
	return &customResourcesUserModel{
		defaultResourcesUserModel: newResourcesUserModel(conn),
	}
}

func (m *customResourcesUserModel) withSession(session sqlx.Session) ResourcesUserModel {
	return NewResourcesUserModel(sqlx.NewSqlConnFromSession(session))
}
