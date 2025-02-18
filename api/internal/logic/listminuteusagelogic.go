package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMinuteUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMinuteUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMinuteUsageLogic {
	return &ListMinuteUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMinuteUsageLogic) ListMinuteUsage(req *types.MinuteUsageListRequest) (resp *types.MinuteUsageListResponse, err error) {
	// Ensure valid pagination
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Prepare optional filters
	var usageId, orgId, userId *uint64
	var startDatetime, endDatetime *string

	if req.UsageId != 0 {
		usageId = &req.UsageId
	}
	if req.OrgId != 0 {
		orgId = &req.OrgId
	}
	if req.UserId != 0 {
		userId = &req.UserId
	}
	if req.StartDatetime != "" {
		startDatetime = &req.StartDatetime
	}
	if req.EndDatetime != "" {
		endDatetime = &req.EndDatetime
	}

	// Get data from database
	items, total, err := l.svcCtx.MinuteUsageModel.List(l.ctx, usageId, orgId, userId, startDatetime, endDatetime, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// Convert model items to response items
	respItems := make([]types.MinuteUsageItem, len(items))
	for i, item := range items {
		respItems[i] = types.MinuteUsageItem{
			Id:            item.Id,
			UsageId:       item.UsageId,
			OrgId:         item.OrgId,
			UserId:        item.UserId,
			ResourceId:    item.ResourceId,
			UsageDatetime: item.UsageDatetime.Format("2006-01-02 15:04:05"),
			Fee:           item.Fee,
			Discount:      item.Discount,
		}
	}

	return &types.MinuteUsageListResponse{
		Response: types.Response{
			Code:    0,
			Message: "success",
		},
		Data: types.MinuteUsageListResponseData{
			Items: respItems,
			Total: total,
		},
	}, nil
}
