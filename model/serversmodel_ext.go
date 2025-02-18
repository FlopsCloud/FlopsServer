package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (m *defaultServersModel) FindList(ctx context.Context, condition string) (*[]Servers, error) {
	query := fmt.Sprintf("select %s from %s %s", strings.Join(serversFieldNames, ","), m.table, condition)
	var resp []Servers
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

func (m *defaultServersModel) Count(ctx context.Context, condition string) (int, error) {
	query := fmt.Sprintf("select count(*) from %s %s", m.table, condition)
	count := 0
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
