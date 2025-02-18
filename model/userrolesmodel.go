package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserRolesModel = (*customUserRolesModel)(nil)

type (
	// UserRolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRolesModel.
	UserRolesModel interface {
		userRolesModel
		WithSession(session sqlx.Session) UserRolesModel
		FindOneByUserIdRoleId(ctx context.Context, userId, roleId uint64) (*UserRoles, error)
		DeleteByRoleId(ctx context.Context, roleId uint64) error
		FindByUserId(ctx context.Context, userId uint64) ([]*UserRoles, error)
	}

	customUserRolesModel struct {
		*defaultUserRolesModel
	}
)

// NewUserRolesModel returns a model for the database table.
func NewUserRolesModel(conn sqlx.SqlConn) UserRolesModel {
	return &customUserRolesModel{
		defaultUserRolesModel: newUserRolesModel(conn),
	}
}

func (m *customUserRolesModel) WithSession(session sqlx.Session) UserRolesModel {
	return NewUserRolesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUserRolesModel) FindOneByUserIdRoleId(ctx context.Context, userId, roleId uint64) (*UserRoles, error) {
	query := fmt.Sprintf("select %s from %s where user_id = ? and role_id = ? limit 1", userRolesRows, m.table)
	var resp UserRoles
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, roleId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserRolesModel) DeleteByRoleId(ctx context.Context, roleId uint64) error {
	query := fmt.Sprintf("delete from %s where role_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, roleId)
	return err
}

func (m *customUserRolesModel) FindByUserId(ctx context.Context, userId uint64) ([]*UserRoles, error) {
	query := fmt.Sprintf("select %s from %s where user_id = ?", userRolesRows, m.table)
	var resp []*UserRoles
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
