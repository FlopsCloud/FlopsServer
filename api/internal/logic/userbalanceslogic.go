package logic

import (
	"context"
	"errors"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBalancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserBalancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBalancesLogic {
	return &UserBalancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserBalancesLogic) Update(userId uint64, money int64, currencyCode string) (*model.Balances, error) {
	item, err := l.svcCtx.BalancesModel.FindOneByUserAndCurrency(l.ctx, userId, currencyCode)
	if err == model.ErrNotFound {
		item = &model.Balances{
			UserId:       userId,
			OrgId:        1, // TODO:关联组织
			Balance:      0,
			CurrencyCode: currencyCode,
		}
	} else if err != nil {
		return item, err
	}

	item.Balance += money
	if item.BalanceId > 0 {
		err = l.svcCtx.BalancesModel.UpdateBalance(l.ctx, item)
	} else {
		_, err = l.svcCtx.BalancesModel.Insert(l.ctx, item)
	}
	return item, err
}

func (l *UserBalancesLogic) Recharge(userId uint64, money int64, orderNo string, payType string, currencyCode string) error {
	moneyFloat := float64(money) / 100
	detail := fmt.Sprintf("%s账户：%s充值%.2f元", currencyCode, payType, moneyFloat)
	err, _ := l.Increase(userId, money, orderNo, payType, currencyCode, detail)
	return err
}

func (l *UserBalancesLogic) Increase(userId uint64, money int64, orderNo string, payType string, currencyCode string, detail string) (error error, balance *model.Balances) {
	item, err := l.Update(userId, money, currencyCode)
	if err != nil {
		return err, nil
	}

	transLogic := NewCreateTransactionRecordsLogic(l.ctx, l.svcCtx)
	resp := transLogic.CreateTransactionRecords(&types.CreateTransactionRecordsRequest{
		TransType:    1, // 1-入金; 2-出金
		CurrencyCode: currencyCode,
		OrgId:        item.OrgId,
		UserId:       userId,
		PayType:      payType,
		Detail:       detail,
		Amount:       money,
		OrderNo:      orderNo,
		Balance:      item.Balance,
	})

	if resp.Code == response.SuccessCode {
		return nil, item
	}
	return errors.New(resp.Message), nil
}

func (l *UserBalancesLogic) Decrease(userId uint64, money int64, orderNo string, currencyCode string, detail string) (error error, balance *model.Balances) {
	item, err := l.Update(userId, -money, currencyCode)
	if err != nil {
		return err, nil
	}

	transLogic := NewCreateTransactionRecordsLogic(l.ctx, l.svcCtx)
	resp := transLogic.CreateTransactionRecords(&types.CreateTransactionRecordsRequest{
		TransType:    2, // 1-入金; 2-出金
		OrgId:        item.OrgId,
		CurrencyCode: currencyCode,
		UserId:       userId,
		PayType:      "",
		Detail:       detail,
		Amount:       money,
		OrderNo:      orderNo,
		Balance:      item.Balance,
	})

	if resp.Code == response.SuccessCode {
		return nil, item
	}
	return errors.New(resp.Message), nil
}

func (l *UserBalancesLogic) ManualAdjust(userId uint64, money int64, currencyCode string, detail string) (error error, balance *model.Balances) {
	item, err := l.Update(userId, money, currencyCode)
	if err != nil {
		return err, nil
	}

	transLogic := NewCreateTransactionRecordsLogic(l.ctx, l.svcCtx)
	resp := transLogic.CreateTransactionRecords(&types.CreateTransactionRecordsRequest{
		TransType:    3, // 1-入金; 2-出金;3-其他
		OrgId:        item.OrgId,
		CurrencyCode: currencyCode,
		UserId:       userId,
		PayType:      "",
		Detail:       detail,
		Amount:       money,
		OrderNo:      "",
		Balance:      item.Balance,
	})

	if resp.Code == response.SuccessCode {
		return nil, item
	}
	return errors.New(resp.Message), nil
}
