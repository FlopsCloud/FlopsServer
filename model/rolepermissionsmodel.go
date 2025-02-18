package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolePermissionsModel = (*customRolePermissionsModel)(nil)

type (
	// RolePermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePermissionsModel.
	RolePermissionsModel interface {
		rolePermissionsModel
		WithSession(session sqlx.Session) RolePermissionsModel
		DeleteByRoleId(ctx context.Context, roleId uint64) error
		DeleteByPermissionId(ctx context.Context, permissionId uint64) error
	}

	customRolePermissionsModel struct {
		*defaultRolePermissionsModel
	}
)

// NewRolePermissionsModel returns a model for the database table.
func NewRolePermissionsModel(conn sqlx.SqlConn) RolePermissionsModel {
	return &customRolePermissionsModel{
		defaultRolePermissionsModel: newRolePermissionsModel(conn),
	}
}

func (m *customRolePermissionsModel) WithSession(session sqlx.Session) RolePermissionsModel {
	return NewRolePermissionsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customRolePermissionsModel) DeleteByRoleId(ctx context.Context, roleId uint64) error {
	query := fmt.Sprintf("delete from %s where role_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *customRolePermissionsModel) DeleteByPermissionId(ctx context.Context, permissionId uint64) error {
	query := fmt.Sprintf("delete from %s where permission_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, permissionId)
	return err
}
