package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RegionsModel = (*customRegionsModel)(nil)

type (
	// RegionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRegionsModel.
	RegionsModel interface {
		regionsModel
		withSession(session sqlx.Session) RegionsModel
		FindAll(ctx context.Context) ([]*Regions, error)
		FindByName(ctx context.Context, name string) (*Regions, error)
	}

	customRegionsModel struct {
		*defaultRegionsModel
	}
)

// NewRegionsModel returns a model for the database table.
func NewRegionsModel(conn sqlx.SqlConn) RegionsModel {
	return &customRegionsModel{
		defaultRegionsModel: newRegionsModel(conn),
	}
}

func (m *customRegionsModel) withSession(session sqlx.Session) RegionsModel {
	return NewRegionsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customRegionsModel) FindAll(ctx context.Context) ([]*Regions, error) {
	var regions []*Regions
	query := fmt.Sprintf("select %s from %s", regionsRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &regions, query)
	return regions, err
}

func (m *customRegionsModel) FindByName(ctx context.Context, name string) (*Regions, error) {
	var region Regions
	query := fmt.Sprintf("select %s from %s where `region_name` = ? limit 1", regionsRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &region, query, name)
	switch err {
	case nil:
		return &region, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
