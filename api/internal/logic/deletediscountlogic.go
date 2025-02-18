package logic

import (
	"context"
	"database/sql"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDiscountLogic {
	return &DeleteDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDiscountLogic) DeleteDiscount(req *types.DeleteDiscountRequest) (resp *types.Response, err error) {
	// Check if discount exists
	existingDiscount, err := l.svcCtx.DiscountsModel.FindOne(l.ctx, req.DiscountId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: "Discount not found",
			}, nil
		}
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to get discount",
			Info:    err.Error(),
		}, nil
	}

	// Soft delete by updating fields
	existingDiscount.IsDeleted = 1
	existingDiscount.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	err = l.svcCtx.DiscountsModel.Update(l.ctx, existingDiscount)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to delete discount",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}, nil
}
