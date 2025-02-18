package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BalancesModel = (*customBalancesModel)(nil)

type (
	// BalancesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBalancesModel.
	BalancesModel interface {
		balancesModel
		WithSession(session sqlx.Session) BalancesModel
		FindList(ctx context.Context, condition string) (*[]Balances, error)
		FindOneByUserAndCurrency(ctx context.Context, userId uint64, currencyCode string) (*Balances, error)
		UpdateBalance(ctx context.Context, data *Balances) error
	}

	customBalancesModel struct {
		*defaultBalancesModel
	}
)

// NewBalancesModel returns a model for the database table.
func NewBalancesModel(conn sqlx.SqlConn) BalancesModel {
	return &customBalancesModel{
		defaultBalancesModel: newBalancesModel(conn),
	}
}

func (m *customBalancesModel) WithSession(session sqlx.Session) BalancesModel {
	return NewBalancesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBalancesModel) FindOneByUserAndCurrency(ctx context.Context, userId uint64, currencyCode string) (*Balances, error) {
	var resp Balances
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `currency_code` = ? limit 1", balancesRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, currencyCode)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customBalancesModel) UpdateBalance(ctx context.Context, data *Balances) error {
	query := fmt.Sprintf("update %s set balance=?,updated_at=? where `user_id` = ? and `currency_code` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Balance, time.Now(), data.UserId, data.CurrencyCode)
	return err
}

func (m *customBalancesModel) FindList(ctx context.Context, condition string) (*[]Balances, error) {
	query := fmt.Sprintf("select %s from %s %s", balancesRows, m.table, condition)
	var resp []Balances
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
