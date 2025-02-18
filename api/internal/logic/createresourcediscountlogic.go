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

type CreateResourceDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建资源折扣
func NewCreateResourceDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateResourceDiscountLogic {
	return &CreateResourceDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateResourceDiscountLogic) CreateResourceDiscount(req *types.CreateResourceDiscountRequest) (resp *types.Response, err error) {
	// Create new resource discount record
	resourceDiscount := model.ResourceDiscounts{
		OrgId:           req.OrgId,
		Memo:            req.Memo,
		HourlyDiscount:  req.HourlyDiscount,
		DailyDiscount:   req.DailyDiscount,
		MonthlyDiscount: req.MonthlyDiscount,
		YearlyDiscount:  req.YearlyDiscount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err = l.svcCtx.ResourceDiscountsModel.Insert(l.ctx, &resourceDiscount)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to create resource discount",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Resource discount created successfully",
	}, nil
}
