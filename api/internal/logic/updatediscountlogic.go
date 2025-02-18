package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDiscountLogic {
	return &UpdateDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDiscountLogic) UpdateDiscount(req *types.UpdateDiscountRequest) (resp *types.DiscountResponse, err error) {
	// Check if discount exists
	existingDiscount, err := l.svcCtx.DiscountsModel.FindOne(l.ctx, req.DiscountId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.DiscountResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Discount not found",
				},
			}, nil
		}
		return &types.DiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get discount",
				Info:    err.Error(),
			},
		}, nil
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.Startdate)
	if err != nil {
		return &types.DiscountResponse{
			Response: types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "Invalid start date format",
				Info:    err.Error(),
			},
		}, nil
	}

	endDate, err := time.Parse("2006-01-02", req.Enddate)
	if err != nil {
		return &types.DiscountResponse{
			Response: types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "Invalid end date format",
				Info:    err.Error(),
			},
		}, nil
	}

	if endDate.Before(startDate) {
		return &types.DiscountResponse{
			Response: types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "End date must be after start date",
			},
		}, nil
	}

	// Update fields
	existingDiscount.Memo.String = req.Memo
	existingDiscount.Memo.Valid = len(req.Memo) > 0
	existingDiscount.Startdate = startDate
	existingDiscount.Enddate = endDate
	existingDiscount.Discount = req.Discount

	err = l.svcCtx.DiscountsModel.Update(l.ctx, existingDiscount)
	if err != nil {
		return &types.DiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to update discount",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get updated discount
	updatedDiscount, err := l.svcCtx.DiscountsModel.FindOne(l.ctx, req.DiscountId)
	if err != nil {
		return &types.DiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Discount updated but failed to retrieve",
				Info:    err.Error(),
			},
		}, nil
	}

	return &types.DiscountResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.Discount{
			DiscountId: updatedDiscount.DiscountId,
			OrgId:      updatedDiscount.OrgId,
			ResourceId: updatedDiscount.ResourceId,
			DiscountBase: types.DiscountBase{
				Memo:      updatedDiscount.Memo.String,
				Startdate: updatedDiscount.Startdate.Format("2006-01-02"),
				Enddate:   updatedDiscount.Enddate.Format("2006-01-02"),
				Discount:  updatedDiscount.Discount,
			},
			IsDeleted: updatedDiscount.IsDeleted,
			DeletedAt: updatedDiscount.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			CreatedAt: updatedDiscount.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: updatedDiscount.CreatedBy,
		},
	}, nil
}
