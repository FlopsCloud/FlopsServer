package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolesModel = (*customRolesModel)(nil)

type (
	// RolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolesModel.
	RolesModel interface {
		rolesModel
		WithSession(session sqlx.Session) RolesModel
		FindByName(ctx context.Context, name string, page, pageSize uint64) ([]*Roles, uint64, error)
	}

	customRolesModel struct {
		*defaultRolesModel
	}
)

// NewRolesModel returns a model for the database table.
func NewRolesModel(conn sqlx.SqlConn) RolesModel {
	return &customRolesModel{
		defaultRolesModel: newRolesModel(conn),
	}
}

func (m *customRolesModel) WithSession(session sqlx.Session) RolesModel {
	return NewRolesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customRolesModel) FindByName(ctx context.Context, name string, page, pageSize uint64) ([]*Roles, uint64, error) {
	where := "1=1"
	var args []interface{}
	if name != "" {
		where += " AND role_name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, where)
	var total uint64
	err := m.conn.QueryRowCtx(ctx, &total, query, args...)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*Roles{}, 0, nil
	}

	query = fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY role_id DESC LIMIT ?,?", rolesRows, m.table, where)
	args = append(args, (page-1)*pageSize, pageSize)
	var resp []*Roles
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
