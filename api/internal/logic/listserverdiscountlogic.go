package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListServerDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListServerDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListServerDiscountLogic {
	return &ListServerDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListServerDiscountLogic) ListServerDiscount(req *types.ListServerDiscountRequest) (resp *types.ListServerDiscountResponse, err error) {
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
	total, err := l.svcCtx.ServerDiscountsModel.Count(l.ctx, req.OrgId, req.ServerId)
	if err != nil {
		return &types.ListServerDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get total count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get server discounts with pagination
	serverDiscounts, err := l.svcCtx.ServerDiscountsModel.FindMany(l.ctx, req.OrgId, req.ServerId, offset, req.PageSize)
	if err != nil {
		return &types.ListServerDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get server discounts",
				Info:    err.Error(),
			},
		}, nil
	}

	// Convert to response type
	var discountList []types.ServerDiscount
	for _, d := range serverDiscounts {
		discountList = append(discountList, types.ServerDiscount{
			SrvDiscountId: d.SrvDiscountId,
			OrgId:         d.OrgId,
			ServerId:      d.ServerId,
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

	return &types.ListServerDiscountResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ListServerDiscountResponseData{
			ServerDiscounts: discountList,
			Total:           total,
		},
	}, nil
}
