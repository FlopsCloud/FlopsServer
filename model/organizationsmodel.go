package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrganizationsModel = (*customOrganizationsModel)(nil)

type (
	// OrganizationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrganizationsModel.
	OrganizationsModel interface {
		organizationsModel
		withSession(session sqlx.Session) OrganizationsModel
		FindByUserId(ctx context.Context, userId uint64) ([]*Organizations, error)
	}

	customOrganizationsModel struct {
		*defaultOrganizationsModel
	}
)

// NewOrganizationsModel returns a model for the database table.
func NewOrganizationsModel(conn sqlx.SqlConn) OrganizationsModel {
	return &customOrganizationsModel{
		defaultOrganizationsModel: newOrganizationsModel(conn),
	}
}

func (m *customOrganizationsModel) withSession(session sqlx.Session) OrganizationsModel {
	return NewOrganizationsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customOrganizationsModel) FindByUserId(ctx context.Context, userId uint64) ([]*Organizations, error) {
	query := fmt.Sprintf("select %s from %s where user_id = ?", organizationsRows, m.table)
	var resp []*Organizations
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
