package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (m *defaultBaseDataModel) FindByType(ctx context.Context, datatype string) (*[]BaseData, error) {
	query := fmt.Sprintf("select %s from %s where `data_type` = ?", baseDataRows, m.table)
	var resp []BaseData
	err := m.conn.QueryRowsCtx(ctx, &resp, query, datatype)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultBaseDataModel) FindByTypeName(ctx context.Context, datatype string, value string) (int, error) {
	query := fmt.Sprintf("select `id` from %s where `data_type` = ? and `name` = ?", m.table)
	count := 0
	err := m.conn.QueryRowCtx(ctx, &count, query, datatype, value)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *defaultBaseDataModel) FindList(ctx context.Context, condition string) (*[]BaseData, error) {
	query := fmt.Sprintf("select %s from %s %s", baseDataRows, m.table, condition)
	var resp []BaseData
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultBaseDataModel) Count(ctx context.Context, condition string) (int, error) {
	query := fmt.Sprintf("select count(1) from %s %s", m.table, condition)
	count := 0
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
