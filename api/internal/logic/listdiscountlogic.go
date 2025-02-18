package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDiscountLogic {
	return &ListDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDiscountLogic) ListDiscount(req *types.ListDiscountRequest) (resp *types.ListDiscountResponse, err error) {
	// Set default page size if not provided
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Get total count
	total, err := l.svcCtx.DiscountsModel.Count(l.ctx, req.OrgId, req.ResourceId)
	if err != nil {
		return &types.ListDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get total count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get discounts with pagination
	discounts, err := l.svcCtx.DiscountsModel.FindMany(l.ctx, req.OrgId, req.ResourceId, offset, req.PageSize)
	if err != nil {
		return &types.ListDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get discounts",
				Info:    err.Error(),
			},
		}, nil
	}

	// Convert to response type
	var discountList []types.Discount
	for _, d := range discounts {
		discountList = append(discountList, types.Discount{
			DiscountId: d.DiscountId,
			OrgId:      d.OrgId,
			ResourceId: d.ResourceId,
			DiscountBase: types.DiscountBase{
				Memo:      d.Memo.String,
				Startdate: d.Startdate.Format("2006-01-02"),
				Enddate:   d.Enddate.Format("2006-01-02"),
				Discount:  d.Discount,
			},
			IsDeleted: d.IsDeleted,
			DeletedAt: d.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			CreatedAt: d.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: d.CreatedBy,
		})
	}

	return &types.ListDiscountResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ListDiscountResponseData{
			Discounts: discountList,
			Total:     total,
		},
	}, nil
}
