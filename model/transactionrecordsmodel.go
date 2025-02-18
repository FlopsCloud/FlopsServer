package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TransactionRecordsModel = (*customTransactionRecordsModel)(nil)

type (
	// TransactionRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTransactionRecordsModel.
	TransactionRecordsModel interface {
		transactionRecordsModel
		FindList(ctx context.Context, condition string) (*[]TransactionRecords, error)
		Count(ctx context.Context, condition string) (int, error)
		withSession(session sqlx.Session) TransactionRecordsModel
	}

	customTransactionRecordsModel struct {
		*defaultTransactionRecordsModel
	}
)

// NewTransactionRecordsModel returns a model for the database table.
func NewTransactionRecordsModel(conn sqlx.SqlConn) TransactionRecordsModel {
	return &customTransactionRecordsModel{
		defaultTransactionRecordsModel: newTransactionRecordsModel(conn),
	}
}

func (m *customTransactionRecordsModel) withSession(session sqlx.Session) TransactionRecordsModel {
	return NewTransactionRecordsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTransactionRecordsModel) FindList(ctx context.Context, condition string) (*[]TransactionRecords, error) {
	query := fmt.Sprintf("select %s from %s %s", transactionRecordsRows, m.table, condition)
	var resp []TransactionRecords
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

func (m *customTransactionRecordsModel) Count(ctx context.Context, condition string) (int, error) {
	query := fmt.Sprintf("select count(1) from %s %s", m.table, condition)
	count := 0
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
