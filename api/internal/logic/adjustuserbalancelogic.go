package logic

import (
	"context"
	"fca/common/response"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdjustUserBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdjustUserBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdjustUserBalanceLogic {
	return &AdjustUserBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdjustUserBalanceLogic) AdjustUserBalance(req *types.AdjustUserBalanceRequest) response.Response {
	userBalanceLogic := NewUserBalancesLogic(l.ctx, l.svcCtx)
	err, item := userBalanceLogic.ManualAdjust(req.UserId, req.Amount, req.CurrencyCode, req.Reason)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	data := make([]types.UserBalanceData, 1)
	data[0] = types.UserBalanceData{
		UserId:       item.UserId,
		Balance:      item.Balance,
		CurrencyCode: item.CurrencyCode,
	}
	return response.OK(data)
}
