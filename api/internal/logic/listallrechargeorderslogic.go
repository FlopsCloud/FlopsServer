package logic

import (
	"context"
	"fca/common/response"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllRechargeOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllRechargeOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllRechargeOrdersLogic {
	return &ListAllRechargeOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllRechargeOrdersLogic) ListAllRechargeOrders(req *types.ListAllRechargeOrdersRequest) response.Response {
	condition := "where 1=1"
	status := 0
	switch req.Status {
	case "pending":
		status = 1
		break
	case "completed":
		status = 2
		break
	case "rejected":
		status = 3
		break
	case "failed":
		status = 4
		break
	}
	if status > 0 {
		condition += fmt.Sprintf(" and status = %d", status)
	}
	if req.PayMethod > 0 {
		condition += fmt.Sprintf(" and pay_method = %d", req.PayMethod)
	}

	total, err := l.svcCtx.RechargeOrderModel.Count(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	condition += fmt.Sprintf(" limit %d,%d", (req.Page-1)*req.PageSize, req.PageSize)

	data, err := l.svcCtx.RechargeOrderModel.FindList(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	count := len(*data)
	orders := make([]types.RechargeOrderData, count)
	for i, item := range *data {
		orders[i] = types.RechargeOrderData{
			UserId:    item.UserId,
			Amount:    item.Amount,
			PayMethod: item.PayMethod,
			OrderNo:   item.OrderNo,
			Status:    l.svcCtx.RechargeOrderModel.GetStatusText(l.ctx, item.Status),
			CreatedAt: item.CreatedAt,
		}
	}

	return response.OK(&types.ListAllRechargeOrdersData{Orders: orders, Total: uint64(total)})
}
