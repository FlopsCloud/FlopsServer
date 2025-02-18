package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RechargeOrdersModel = (*customRechargeOrdersModel)(nil)

type (
	// RechargeOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRechargeOrdersModel.
	RechargeOrdersModel interface {
		rechargeOrdersModel
		FindList(ctx context.Context, condition string) (*[]RechargeOrders, error)
		Count(ctx context.Context, condition string) (int, error)
		FindOneByOrderId(ctx context.Context, orderNo string) (*RechargeOrders, error)
		GetStatusText(ctx context.Context, status uint64) string
		withSession(session sqlx.Session) RechargeOrdersModel
	}

	customRechargeOrdersModel struct {
		*defaultRechargeOrdersModel
	}
)

// NewRechargeOrdersModel returns a model for the database table.
func NewRechargeOrdersModel(conn sqlx.SqlConn) RechargeOrdersModel {
	return &customRechargeOrdersModel{
		defaultRechargeOrdersModel: newRechargeOrdersModel(conn),
	}
}

func (m *customRechargeOrdersModel) withSession(session sqlx.Session) RechargeOrdersModel {
	return NewRechargeOrdersModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customRechargeOrdersModel) FindOneByOrderId(ctx context.Context, orderNo string) (*RechargeOrders, error) {
	var resp RechargeOrders
	query := fmt.Sprintf("select %s from %s where `order_no` = ? limit 1", rechargeOrdersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderNo)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customRechargeOrdersModel) FindList(ctx context.Context, condition string) (*[]RechargeOrders, error) {
	query := fmt.Sprintf("select %s from %s %s", rechargeOrdersRows, m.table, condition)
	var resp []RechargeOrders
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

func (m *customRechargeOrdersModel) Count(ctx context.Context, condition string) (int, error) {
	query := fmt.Sprintf("select count(1) from %s %s", m.table, condition)
	count := 0
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *customRechargeOrdersModel) GetStatusText(ctx context.Context, status uint64) string {
	result := "undefined"
	switch status {
	case 2:
		result = "completed"
		break
	case 1:
		result = "pending"
		break
	case 3:
		result = "rejected"
		break
	case 4:
		result = "failed"
		break
	}
	return result
}
