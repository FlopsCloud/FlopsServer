package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiscountLogic {
	return &GetDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDiscountLogic) GetDiscount(req *types.GetDiscountRequest) (resp *types.DiscountResponse, err error) {
	discount, err := l.svcCtx.DiscountsModel.FindOne(l.ctx, req.DiscountId)
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

	return &types.DiscountResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.Discount{
			DiscountId: discount.DiscountId,
			OrgId:      discount.OrgId,
			ResourceId: discount.ResourceId,
			DiscountBase: types.DiscountBase{
				Memo:      discount.Memo.String,
				Startdate: discount.Startdate.Format("2006-01-02"),
				Enddate:   discount.Enddate.Format("2006-01-02"),
				Discount:  discount.Discount,
			},
			IsDeleted: discount.IsDeleted,
			DeletedAt: discount.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			CreatedAt: discount.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: discount.CreatedBy,
		},
	}, nil
}
