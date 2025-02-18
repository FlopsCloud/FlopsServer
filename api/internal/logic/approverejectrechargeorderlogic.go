package logic

import (
	"context"
	"fca/common/response"
	"fmt"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveRejectRechargeOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveRejectRechargeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveRejectRechargeOrderLogic {
	return &ApproveRejectRechargeOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveRejectRechargeOrderLogic) ApproveRejectRechargeOrder(req *types.ApproveRejectRechargeOrderRequest) response.Response {
	rechargeModel := l.svcCtx.RechargeOrderModel
	order, err := rechargeModel.FindOneByOrderId(l.ctx, req.OrderNo)
	if err != nil {
		return response.Fail(response.RecordNotExistCode, err.Error())
	}
	if order.Status != 1 {
		statusText := rechargeModel.GetStatusText(l.ctx, order.Status)
		return response.Fail(response.ServerErrorCode, fmt.Sprintf("order status is %s", statusText))
	}

	if req.Action == "approve" {
		userLogic := NewUserBalancesLogic(l.ctx, l.svcCtx)
		payType := "其他"
		switch order.PayMethod {
		case 1:
			payType = "微信支付"
			break
		case 2:
			payType = "支付宝"
			break
		}
		err, _ = userLogic.Increase(order.UserId, order.Amount, order.OrderNo, payType, "CNY", req.Reason)
		if err == nil {
			order.Status = 2
			order.PaidAt = uint64(time.Now().UnixMicro())
			order.Remark = req.Reason
			err = rechargeModel.Update(l.ctx, order)
		}
	} else if req.Action == "reject" {
		order.Status = 3
		order.Remark = req.Reason
		err = rechargeModel.Update(l.ctx, order)
	} else {
		return response.Fail(response.ServerErrorCode, fmt.Sprintf("action '%s' is invaild", req.Action))
	}
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	return response.OK(order.OrderNo)
}
