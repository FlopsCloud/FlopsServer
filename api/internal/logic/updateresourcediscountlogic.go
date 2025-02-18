package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateResourceDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新资源折扣
func NewUpdateResourceDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateResourceDiscountLogic {
	return &UpdateResourceDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateResourceDiscountLogic) UpdateResourceDiscount(req *types.UpdateResourceDiscountRequest) (resp *types.Response, err error) {
	// Find existing resource discount
	existingDiscount, err := l.svcCtx.ResourceDiscountsModel.FindOne(l.ctx, req.DiscountId)
	if err != nil {
		return &types.Response{
			Code:    response.NotFoundCode,
			Message: "Resource discount not found",
			Info:    err.Error(),
		}, nil
	}

	// Update fields
	existingDiscount.OrgId = req.OrgId
	existingDiscount.Memo = req.Memo
	existingDiscount.HourlyDiscount = req.HourlyDiscount
	existingDiscount.DailyDiscount = req.DailyDiscount
	existingDiscount.MonthlyDiscount = req.MonthlyDiscount
	existingDiscount.YearlyDiscount = req.YearlyDiscount
	existingDiscount.UpdatedAt = time.Now()

	err = l.svcCtx.ResourceDiscountsModel.Update(l.ctx, existingDiscount)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to update resource discount",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Resource discount updated successfully",
	}, nil
}
