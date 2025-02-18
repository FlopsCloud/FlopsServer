package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadResourceDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 读取资源折扣
func NewReadResourceDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadResourceDiscountLogic {
	return &ReadResourceDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadResourceDiscountLogic) ReadResourceDiscount(req *types.ReadResourceDiscountRequest) (resp *types.ReadResourceDiscountResponse, err error) {
	discount, err := l.svcCtx.ResourceDiscountsModel.FindOne(l.ctx, req.DiscountId)
	if err != nil {
		return &types.ReadResourceDiscountResponse{
			Response: types.Response{
				Code:    response.NotFoundCode,
				Message: "Resource discount not found",
				Info:    err.Error(),
			},
		}, nil
	}

	return &types.ReadResourceDiscountResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ResourceDiscount{
			DiscountId:      discount.DiscountId,
			OrgId:           discount.OrgId,
			Memo:            discount.Memo,
			HourlyDiscount:  discount.HourlyDiscount,
			DailyDiscount:   discount.DailyDiscount,
			MonthlyDiscount: discount.MonthlyDiscount,
			YearlyDiscount:  discount.YearlyDiscount,
			CreatedAt:       uint64(discount.CreatedAt.Unix()),
			UpdatedAt:       uint64(discount.UpdatedAt.Unix()),
		},
	}, nil
}
