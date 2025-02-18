package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListResourceDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 读取资源折扣
func NewListResourceDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListResourceDiscountLogic {
	return &ListResourceDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListResourceDiscountLogic) ListResourceDiscount(req *types.ListResourceDiscountRequest) (resp *types.ListResourceDiscountResponse, err error) {
	// Initialize response
	resp = &types.ListResourceDiscountResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ListResourceDiscountResponseData{
			Discounts: make([]types.ResourceDiscount, 0),
			Total:     0,
		},
	}

	// Set pagination defaults if not provided
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// Query total count
	total, err := l.svcCtx.ResourceDiscountsModel.Count(l.ctx, req.OrgId)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "Server Error "
		resp.Info = "Count:" + err.Error()
		return
	}

	// Query discounts with pagination
	discounts, err := l.svcCtx.ResourceDiscountsModel.FindByConditions(
		l.ctx,
		req.OrgId,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "Server Error "
		resp.Info = "FindByConditions:" + err.Error()
		return
	}

	// Convert model discounts to response type
	for _, discount := range discounts {
		resp.Data.Discounts = append(resp.Data.Discounts, types.ResourceDiscount{
			DiscountId:      discount.DiscountId,
			OrgId:           discount.OrgId,
			Memo:            discount.Memo,
			HourlyDiscount:  discount.HourlyDiscount,
			DailyDiscount:   discount.DailyDiscount,
			MonthlyDiscount: discount.MonthlyDiscount,
			YearlyDiscount:  discount.YearlyDiscount,
			CreatedAt:       uint64(discount.CreatedAt.Unix()),
			UpdatedAt:       uint64(discount.UpdatedAt.Unix()),
		})
	}

	// Set total count
	resp.Data.Total = uint64(total)

	return resp, nil

}
