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

type DeleteServerDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteServerDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteServerDiscountLogic {
	return &DeleteServerDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteServerDiscountLogic) DeleteServerDiscount(req *types.DeleteServerDiscountRequest) (resp *types.Response, err error) {
	// Check if server discount exists
	existingDiscount, err := l.svcCtx.ServerDiscountsModel.FindOne(l.ctx, req.SrvDiscountId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: "Server discount not found",
			}, nil
		}
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to get server discount",
			Info:    err.Error(),
		}, nil
	}

	// Soft delete by updating fields
	existingDiscount.IsDeleted = 1
	existingDiscount.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	err = l.svcCtx.ServerDiscountsModel.Update(l.ctx, existingDiscount)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to delete server discount",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}, nil
}
