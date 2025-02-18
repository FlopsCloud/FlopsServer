package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListHourlyUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListHourlyUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListHourlyUsageLogic {
	return &AdminListHourlyUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListHourlyUsageLogic) AdminListHourlyUsage(req *types.AdminHourlyUsageListRequest) (resp *types.HourlyUsageListResponse, err error) {
	role, _ := l.ctx.Value("role").(string)

	if role != "superadmin" || role == "admin" {
		return &types.HourlyUsageListResponse{
			Response: types.Response{
				Code:    response.UnauthorizedCode,
				Message: "permission denied",
			},
		}, nil
	}
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

	if req.UserId != 0 {
		userId = &req.UserId
	}

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
			Id:            item.UsageId,
			UsageId:       item.UsageId,
			OrgId:         item.OrgId,
			UserId:        item.UserId,
			RunresId:      item.RunresId,
			UsageDate:     uint64(item.UsageDate.Unix()),
			Fee:           item.Fee,
			Discount:      item.Discount,
			InstanceId:    item.InstanceId,
			Type:          item.Type,
			InstanceName:  item.InstanceName,
			ResourceName:  item.ResourceName,
			Daynum:        item.Daynum,
			Hournum:       item.Hournum,
			MinuteBegin:   item.MinuteBegin,
			MinuteEnd:     item.MinuteEnd,
			MinuteTotal:   item.MinuteTotal,
			IsCharged:     item.IsCharged,
			UnitHourPrice: item.UnitHourPrice,
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
