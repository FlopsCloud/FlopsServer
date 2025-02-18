package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListResourcesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListResourcesLogic {
	return &ListResourcesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListResourcesLogic) ListResources(req *types.ResourceListRequest) (resp *types.ResourceListResponse, err error) {
	// Ensure valid pagination
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Prepare optional filters
	var orgId *uint64
	var resourceType *string

	if req.OrgId != 0 {
		orgId = &req.OrgId
	}
	if req.ResourceType != "" {
		resourceType = &req.ResourceType
	}

	// Get data from database
	items, total, err := l.svcCtx.ResourcesModel.List(l.ctx, orgId, resourceType, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// Convert model items to response items
	respItems := make([]types.ResourceItem, len(items))
	for i, item := range items {
		respItems[i] = types.ResourceItem{
			ResourceId: item.ResourceId,

			ResourceType:  item.ResourceType,
			UnitMinPrice:  item.UnitMinPrice,
			UnitHourPrice: item.UnitHourPrice,
			UnitDayPrice:  item.UnitDayPrice,
			IsDeleted:     item.IsDeleted,
			CreatedAt:     item.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     item.CreatedBy,
		}
	}

	return &types.ResourceListResponse{
		Response: types.Response{
			Code:    0,
			Message: "success",
		},
		Data: types.ResourceListResponseData{
			Items: respItems,
			Total: total,
		},
	}, nil
}
