package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRegionLogic {
	return &AddRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRegionLogic) AddRegion(req *types.AddRegionRequest) (resp *types.Response, err error) {
	// Check if region name already exists
	exists, err := l.svcCtx.RegionsModel.FindByName(l.ctx, req.RegionName)
	if err != nil && err != model.ErrNotFound {
		return &types.Response{
			Code:    500,
			Message: err.Error(),
		}, nil
	}
	if exists != nil {
		return &types.Response{
			Code:    400,
			Message: "Region name already exists",
		}, nil
	}

	// Create new region
	_, err = l.svcCtx.RegionsModel.Insert(l.ctx, &model.Regions{
		RegionName: req.RegionName,
		RegionCode: req.RegionCode,
	})
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "Region added successfully",
	}, nil
}
