package logic

import (
	"context"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListOrderRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListOrderRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListOrderRecordsLogic {
	return &AdminListOrderRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListOrderRecordsLogic) AdminListOrderRecords(req *types.AdminOrderRecordsListRequest) (resp *types.OrderRecordsListResponse, err error) {
	role, _ := l.ctx.Value("role").(string)
	if role != "admin" && role != "superadmin" {
		return &types.OrderRecordsListResponse{
			Response: types.Response{
				Code:    response.UnauthorizedCode,
				Message: "You are not authorized to access this resource.",
			},
		}, nil
	}
	condition := "where 1=1"
	if req.UserId > 0 {
		condition += fmt.Sprintf(" and user_id = %d", req.UserId)
	}
	if req.OrderType > 0 {
		condition += fmt.Sprintf(" and order_type = %d", req.OrderType)
	}
	if req.OrderNo != "" {
		condition += fmt.Sprintf(" and order_no = '%s'", req.OrderNo)
	}
	if req.BeginCreatedAt > 0 {
		condition += fmt.Sprintf(" and created_at > %d", req.BeginCreatedAt)
	}
	if req.EndCreatedAt > 0 {
		condition += fmt.Sprintf(" and created_at < %d", req.EndCreatedAt)
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}

	total, err := l.svcCtx.OrderRecordsModel.Count(l.ctx, condition)
	if err != nil {
		return &types.OrderRecordsListResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: err.Error(),
			},
		}, nil
	}

	condition += fmt.Sprintf(" limit %d,%d", (req.Page-1)*req.PageSize, req.PageSize)

	data, err := l.svcCtx.OrderRecordsModel.FindList(l.ctx, condition)
	if err != nil {
		return &types.OrderRecordsListResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: err.Error(),
			},
		}, nil
	}

	// bytes, err := json.Marshal(data)
	// if err != nil {
	// 	return &types.OrderRecordsListResponse{
	// 		Response: types.Response{
	// 			Code:    response.ServerErrorCode,
	// 			Message: err.Error(),
	// 		},
	// 	}, nil
	// }

	// var items []types.TransactionRecords
	// err = json.Unmarshal(bytes, &items)
	// if err != nil {
	// 	return &types.OrderRecordsListResponse{
	// 		Response: types.Response{
	// 			Code:    response.ServerErrorCode,
	// 			Message: err.Error(),
	// 		},
	// 	}, nil
	// }

	var orders []types.OrderRecords

	// Id            uint64 `json:"id"`
	// OrderNo       string `json:"orderNo"`
	// InstanceId    uint64 `json:"instanceId"`
	// RunningTime   uint64 `json:"runningTime"`
	// CreatedAt     uint64 `json:"createdAt"`
	// BillingMethod uint64 `json:"billingMethod"`
	// OrderAmount   uint64 `json:"orderAmount"`
	// Discount      uint64 `json:"discount"`
	// ActualAmount  uint64 `json:"actualAmount"`
	// UserId        uint64 `json:"userId"`
	// OrderType     uint64 `json:"orderType"`

	// Id            uint64 `db:"id"`             // ID
	// OrderNo       string `db:"order_no"`       // 订单号
	// InstanceId    uint64 `db:"instance_id"`    // 实例ID
	// RunningTime   uint64 `db:"running_time"`   // 运行时长
	// CreatedAt     uint64 `db:"created_at"`     // 创建时间
	// BillingMethod uint64 `db:"billing_method"` // 计费方式（1-按时、2-按天、3-按周、4-按月、0-按量）
	// OrderAmount   uint64 `db:"order_amount"`   // 订单金额
	// Discount      uint64 `db:"discount"`       // 补贴减免
	// ActualAmount  uint64 `db:"actual_amount"`  // 实收金额
	// UserId        uint64 `db:"user_id"`        // 用户ID
	// OrderType     uint64 `db:"order_type"`     // 订单类型（1-收入、2-消费）

	for _, item := range *data {
		order := types.OrderRecords{
			Id:            item.Id,
			InstanceId:    item.InstanceId,
			OrderNo:       item.OrderNo,
			OrderType:     item.OrderType,
			RunningTime:   item.RunningTime,
			BillingMethod: item.BillingMethod,
			OrderAmount:   item.OrderAmount,
			Discount:      item.Discount,
			ActualAmount:  item.ActualAmount,
			UserId:        item.UserId,
			CreatedAt:     uint64(item.CreatedAt.Unix()),
		}
		orders = append(orders, order)
	}
	return &types.OrderRecordsListResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: types.OrderRecordsListResponseData{
			Orders: orders,
			Total:  total,
		},
	}, nil
}
