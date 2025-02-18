package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateResourceLogic {
	return &UpdateResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateResourceLogic) UpdateResource(req *types.ResourceUpdateRequest) (resp *types.ResourceResponse, err error) {
	// Get existing resource
	resource, err := l.svcCtx.ResourcesModel.FindOne(l.ctx, req.ResourceId)
	if err != nil {
		return nil, err
	}

	// Update fields if provided

	if req.ResourceType != "" {
		resource.ResourceType = req.ResourceType
	}
	if req.UnitMinPrice != 0 {
		resource.UnitMinPrice = req.UnitMinPrice
	}
	if req.UnitHourPrice != 0 {
		resource.UnitHourPrice = req.UnitHourPrice
	}
	if req.UnitDayPrice != 0 {
		resource.UnitDayPrice = req.UnitDayPrice
	}

	// Update in database
	err = l.svcCtx.ResourcesModel.Update(l.ctx, resource)
	if err != nil {
		return nil, err
	}

	// Get updated resource
	updatedResource, err := l.svcCtx.ResourcesModel.FindOne(l.ctx, req.ResourceId)
	if err != nil {
		return nil, err
	}

	return &types.ResourceResponse{
		Response: types.Response{
			Code:    0,
			Message: "success",
		},
		Data: types.ResourceItem{
			ResourceId: updatedResource.ResourceId,

			ResourceType:  updatedResource.ResourceType,
			UnitMinPrice:  updatedResource.UnitMinPrice,
			UnitHourPrice: updatedResource.UnitHourPrice,
			UnitDayPrice:  updatedResource.UnitDayPrice,
			IsDeleted:     updatedResource.IsDeleted,
			CreatedAt:     updatedResource.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     updatedResource.CreatedBy,
		},
	}, nil
}
