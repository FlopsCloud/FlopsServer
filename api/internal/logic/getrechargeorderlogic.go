package logic

import (
	"context"
	"fca/common/response"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRechargeOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRechargeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRechargeOrderLogic {
	return &GetRechargeOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRechargeOrderLogic) GetRechargeOrder(req *types.GetRechargeOrderRequest) response.Response {
	order, err := l.svcCtx.RechargeOrderModel.FindOneByOrderId(l.ctx, req.OrderNo)
	if err != nil {
		return response.FailWithInfo(response.ServerErrorCode, "获取充值订单失败", err.Error())
	}
	statusText := l.svcCtx.RechargeOrderModel.GetStatusText(l.ctx, order.Status)
	data := types.RechargeOrderData{
		OrderNo:   order.OrderNo,
		UserId:    order.UserId,
		Amount:    order.Amount,
		CreatedAt: uint64(order.CreatedAt.Unix()),
		PayMethod: order.PayMethod,
		Status:    statusText,
	}
	return response.OK(data)
}
