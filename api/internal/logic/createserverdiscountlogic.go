package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateServerDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateServerDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateServerDiscountLogic {
	return &CreateServerDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateServerDiscountLogic) CreateServerDiscount(req *types.CreateServerDiscountRequest) (resp *types.ServerDiscountResponse, err error) {
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
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	serverDiscount := &model.ServerDiscounts{
		OrgId:    req.OrgId,
		ServerId: req.ServerId,
		Memo: sql.NullString{
			String: req.Memo,
			Valid:  len(req.Memo) > 0,
		},
		Startdate: startDate,
		Enddate:   endDate,
		Discount:  req.Discount,
		IsDeleted: 0,
		CreatedBy: uint64(uid),
	}

	result, err := l.svcCtx.ServerDiscountsModel.Insert(l.ctx, serverDiscount)
	if err != nil {
		return &types.ServerDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to create server discount",
				Info:    err.Error(),
			},
		}, nil
	}

	srvDiscountId, _ := result.LastInsertId()
	createdDiscount, err := l.svcCtx.ServerDiscountsModel.FindOne(l.ctx, uint64(srvDiscountId))
	if err != nil {
		return &types.ServerDiscountResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Server discount created but failed to retrieve",
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
			SrvDiscountId: createdDiscount.SrvDiscountId,
			OrgId:         createdDiscount.OrgId,
			ServerId:      createdDiscount.ServerId,
			DiscountBase: types.DiscountBase{
				Memo:      createdDiscount.Memo.String,
				Startdate: createdDiscount.Startdate.Format("2006-01-02"),
				Enddate:   createdDiscount.Enddate.Format("2006-01-02"),
				Discount:  createdDiscount.Discount,
			},
			IsDeleted: createdDiscount.IsDeleted,
			DeletedAt: createdDiscount.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			CreatedAt: createdDiscount.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: createdDiscount.CreatedBy,
		},
	}, nil
}
