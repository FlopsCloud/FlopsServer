package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteResourceDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除资源折扣
func NewDeleteResourceDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteResourceDiscountLogic {
	return &DeleteResourceDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteResourceDiscountLogic) DeleteResourceDiscount(req *types.DeleteResourceDiscountRequest) (resp *types.Response, err error) {
	err = l.svcCtx.ResourceDiscountsModel.Delete(l.ctx, req.DiscountId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to delete resource discount",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Resource discount deleted successfully",
	}, nil
}
