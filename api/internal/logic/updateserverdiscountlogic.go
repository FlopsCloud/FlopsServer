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

type UpdateServerDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateServerDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateServerDiscountLogic {
	return &UpdateServerDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateServerDiscountLogic) UpdateServerDiscount(req *types.UpdateServerDiscountRequest) (resp *types.ServerDiscountResponse, err error) {
	// Check if server discount exists
	existingDiscount, err := l.svcCtx.ServerDiscountsModel.FindOne(l.ctx, req.SrvDiscountId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.ServerDiscountResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Server discount not found",
				},
			}, nil
		}
		return &types.ServerDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get server discount",
				Info:    err.Error(),
			},
		}, nil
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.Startdate)
	if err != nil {
		return &types.ServerDiscountResponse{
			Response: types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "Invalid start date format",
				Info:    err.Error(),
			},
		}, nil
	}

	endDate, err := time.Parse("2006-01-02", req.Enddate)
	if err != nil {
		return &types.ServerDiscountResponse{
			Response: types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "Invalid end date format",
				Info:    err.Error(),
			},
		}, nil
	}

	if endDate.Before(startDate) {
		return &types.ServerDiscountResponse{
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

	err = l.svcCtx.ServerDiscountsModel.Update(l.ctx, existingDiscount)
	if err != nil {
		return &types.ServerDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to update server discount",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get updated server discount
	updatedDiscount, err := l.svcCtx.ServerDiscountsModel.FindOne(l.ctx, req.SrvDiscountId)
	if err != nil {
		return &types.ServerDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Server discount updated but failed to retrieve",
				Info:    err.Error(),
			},
		}, nil
	}

	return &types.ServerDiscountResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ServerDiscount{
			SrvDiscountId: updatedDiscount.SrvDiscountId,
			OrgId:         updatedDiscount.OrgId,
			ServerId:      updatedDiscount.ServerId,
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
