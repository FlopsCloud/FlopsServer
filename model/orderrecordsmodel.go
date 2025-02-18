package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderRecordsModel = (*customOrderRecordsModel)(nil)

type (
	// OrderRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderRecordsModel.
	OrderRecordsModel interface {
		orderRecordsModel
		FindList(ctx context.Context, condition string) (*[]OrderRecords, error)
		Count(ctx context.Context, condition string) (int, error)
		withSession(session sqlx.Session) OrderRecordsModel
	}

	customOrderRecordsModel struct {
		*defaultOrderRecordsModel
	}
)

// NewOrderRecordsModel returns a model for the database table.
func NewOrderRecordsModel(conn sqlx.SqlConn) OrderRecordsModel {
	return &customOrderRecordsModel{
		defaultOrderRecordsModel: newOrderRecordsModel(conn),
	}
}

func (m *customOrderRecordsModel) withSession(session sqlx.Session) OrderRecordsModel {
	return NewOrderRecordsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customOrderRecordsModel) FindList(ctx context.Context, condition string) (*[]OrderRecords, error) {
	query := fmt.Sprintf("select %s from %s %s", orderRecordsRows, m.table, condition)
	var resp []OrderRecords
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

func (m *customOrderRecordsModel) Count(ctx context.Context, condition string) (int, error) {
	query := fmt.Sprintf("select count(1) from %s %s", m.table, condition)
	count := 0
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
