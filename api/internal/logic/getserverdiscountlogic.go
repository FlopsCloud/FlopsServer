package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServerDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerDiscountLogic {
	return &GetServerDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerDiscountLogic) GetServerDiscount(req *types.GetServerDiscountRequest) (resp *types.ServerDiscountResponse, err error) {
	serverDiscount, err := l.svcCtx.ServerDiscountsModel.FindOne(l.ctx, req.SrvDiscountId)
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

	return &types.ServerDiscountResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ServerDiscount{
			SrvDiscountId: serverDiscount.SrvDiscountId,
			OrgId:         serverDiscount.OrgId,
			ServerId:      serverDiscount.ServerId,
			DiscountBase: types.DiscountBase{
				Memo:      serverDiscount.Memo.String,
				Startdate: serverDiscount.Startdate.Format("2006-01-02"),
				Enddate:   serverDiscount.Enddate.Format("2006-01-02"),
				Discount:  serverDiscount.Discount,
			},
			IsDeleted: serverDiscount.IsDeleted,
			DeletedAt: serverDiscount.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			CreatedAt: serverDiscount.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: serverDiscount.CreatedBy,
		},
	}, nil
}
