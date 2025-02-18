package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRegionLogic {
	return &UpdateRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRegionLogic) UpdateRegion(req *types.UpdateRegionRequest) (resp *types.Response, err error) {
	sysrole, _ := l.ctx.Value("role").(string)
	if sysrole != "admin" && sysrole != "superadmin" {

		return &types.Response{
			Code:    response.UnauthorizedCode,
			Message: "only admin can access",
		}, nil
	}

	// Check if region exists
	Regions, err := l.svcCtx.RegionsModel.FindOne(l.ctx, req.RegionId)
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

	Regions.RegionName = req.RegionName
	Regions.RegionCode = req.RegionCode

	err = l.svcCtx.RegionsModel.Update(l.ctx, Regions)
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "Region updated successfully",
	}, nil
}
