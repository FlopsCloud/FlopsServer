package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListHourlyUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListHourlyUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListHourlyUsageLogic {
	return &ListHourlyUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListHourlyUsageLogic) ListHourlyUsage(req *types.HourlyUsageListRequest) (resp *types.HourlyUsageListResponse, err error) {
	// Ensure valid pagination
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Prepare optional filters
	var usageId, orgId, userId *uint64
	var startDatetime, endDatetime *uint64

	if req.UsageId != 0 {
		usageId = &req.UsageId
	}
	if req.OrgId != 0 {
		orgId = &req.OrgId
	}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	uid64 := uint64(uid)

	userId = &uid64

	if req.StartDatetime != 0 {
		startDatetime = &req.StartDatetime
	}
	if req.EndDatetime != 0 {
		endDatetime = &req.EndDatetime
	}

	// Get data from database
	items, total, err := l.svcCtx.HourlyUsageModel.List(l.ctx, usageId, orgId, userId, startDatetime, endDatetime, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// Convert model items to response items
	respItems := make([]types.HourlyUsageItem, len(items))
	for i, item := range items {
		respItems[i] = types.HourlyUsageItem{
			Id:           item.UsageId,
			UsageId:      item.UsageId,
			OrgId:        item.OrgId,
			UserId:       item.UserId,
			RunresId:     item.RunresId,
			UsageDate:    uint64(item.UsageDate.Unix()),
			Fee:          item.Fee,
			Discount:     item.Discount,
			InstanceId:   item.InstanceId,
			Type:         item.Type,
			InstanceName: item.InstanceName,
			ResourceName: item.ResourceName,
			Daynum:       item.Daynum,
			Hournum:      item.Hournum,
			MinuteBegin:  item.MinuteBegin,
			MinuteEnd:    item.MinuteEnd,
			MinuteTotal:  item.MinuteTotal,
			IsCharged:    item.IsCharged,
		}
	}

	return &types.HourlyUsageListResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: types.HourlyUsageListResponseData{
			Items: respItems,
			Total: total,
		},
	}, nil
}
