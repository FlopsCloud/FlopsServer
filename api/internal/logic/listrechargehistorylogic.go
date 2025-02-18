package logic

import (
	"context"
	"fca/common/response"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRechargeHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRechargeHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRechargeHistoryLogic {
	return &ListRechargeHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRechargeHistoryLogic) ListRechargeHistory(req *types.ListRechargeHistoryRequest) response.Response {
	condition := "where trans_type = 1"
	if req.UserId > 0 {
		condition += fmt.Sprintf(" and user_id = %d", req.UserId)
	}

	total, err := l.svcCtx.TransactionRecordsModel.Count(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	condition += fmt.Sprintf(" limit %d,%d", (req.Page-1)*req.PageSize, req.PageSize)

	data, err := l.svcCtx.TransactionRecordsModel.FindList(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	var items = *data

	count := len(items)
	orders := make([]types.RechargeOrderData, count)
	for i, item := range items {
		orders[i] = types.RechargeOrderData{
			UserId: item.UserId,
			Amount: item.Amount,
			//PayMethod: item.PayType,
			OrderNo:   item.OrderNo,
			Status:    "completed",
			CreatedAt: item.CreatedAt,
		}
	}

	return response.OK(&types.ListRechargeHistoryData{Orders: orders, Total: uint64(total)})
}
