package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRegionLogic {
	return &DeleteRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRegionLogic) DeleteRegion(req *types.DeleteRegionRequest) (resp *types.Response, err error) {
	sysrole, _ := l.ctx.Value("role").(string)
	if sysrole != "admin" && sysrole != "superadmin" {
		var res types.Response
		res.Code = response.UnauthorizedCode
		res.Message = "only admin can access"
		return &res, nil
	}
	// Check if region exists
	_, err = l.svcCtx.RegionsModel.FindOne(l.ctx, req.RegionId)
	if err == model.ErrNotFound {
		return &types.Response{
			Code:    404,
			Message: "Region not found",
		}, nil
	}
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// Delete region
	err = l.svcCtx.RegionsModel.Delete(l.ctx, req.RegionId)
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "Region deleted successfully",
	}, nil
}
