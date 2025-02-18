package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRechargeOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRechargeOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRechargeOrdersLogic {
	return &ListRechargeOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRechargeOrdersLogic) ListRechargeOrders(req *types.ListAllRechargeOrdersRequest) (resp *types.ListAllRechargeOrdersResponse, err error) {

	resp = &types.ListAllRechargeOrdersResponse{}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	// email, _ := l.ctx.Value("email").(string)

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
	condition += fmt.Sprintf(" and user_id = %d", uid)
	if status > 0 {
		condition += fmt.Sprintf(" and status = %d", status)
	}
	if req.PayMethod > 0 {
		condition += fmt.Sprintf(" and pay_method = %d", req.PayMethod)
	}

	condition += " and is_deleted = 0"

	total, err := l.svcCtx.RechargeOrderModel.Count(l.ctx, condition)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "查询订单数量出错"
		resp.Info = err.Error()
		return resp, nil
	}
	condition += fmt.Sprintf(" limit %d,%d", (req.Page-1)*req.PageSize, req.PageSize)

	data, err := l.svcCtx.RechargeOrderModel.FindList(l.ctx, condition)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "查询订单数量出错"
		resp.Info = err.Error()
		return resp, nil

	}

	// Id            int64     `json:"id"`             // ID
	// UserId        uint64    `json:"userId"`        // 用户ID
	// OrgId         uint64    `json:"orgId"`         // 组织ID
	// OrderNo   string `json:"orderNo"`
	// OrderTitle    string    `json:"orderTitle"`    // 订单标题
	// PayMethod uint64 `json:"payMethod"`
	// Remark        string    `json:"remark"`         // 订单备注
	// Amount    int64  `json:"amount"`
	// Status    string `json:"status"`
	// CreatedAt uint64 `json:"createdAt"`
	// PaidAt        uint64    `db:"paidAt"`        // 支付时间
	// PayCodeUrl    string    `db:"payCodeUrl"`   // 付款链接
	// TransactionId string    `db:"transactionId"` // 交易流水号
	count := len(*data)
	orders := make([]types.RechargeOrderData, count)
	for i, item := range *data {

		orders[i] = types.RechargeOrderData{
			Id:            item.Id,
			UserId:        item.UserId,
			OrgId:         item.OrgId,
			OrderNo:       item.OrderNo,
			OrderTitle:    item.OrderTitle,
			PayMethod:     item.PayMethod,
			Remark:        item.Remark,
			Amount:        item.Amount,
			Status:        l.svcCtx.RechargeOrderModel.GetStatusText(l.ctx, item.Status),
			PaidAt:        item.PaidAt,
			TransactionId: item.TransactionId,
			PayCodeUrl:    item.PayCodeUrl,
			CreatedAt:     uint64(item.CreatedAt.Unix()),
		}
	}

	resp.Code = response.SuccessCode
	resp.Message = "查询订单列表成功"
	resp.Data = &types.ListAllRechargeOrdersData{Orders: orders, Total: uint64(total)}
	return resp, nil

}
