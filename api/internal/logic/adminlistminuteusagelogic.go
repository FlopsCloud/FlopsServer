package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListMinuteUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取每分钟使用量 根据 usageid + orgid + user_id + reuresid 或者 instanceid 进行与每分钟关联，金额单位分
func NewAdminListMinuteUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListMinuteUsageLogic {
	return &AdminListMinuteUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListMinuteUsageLogic) AdminListMinuteUsage(req *types.AdminMinuteUsageListRequest) (resp *types.MinuteUsageListResponse, err error) {
	role, _ := l.ctx.Value("role").(string)

	if role != "superadmin" || role == "admin" {
		return &types.MinuteUsageListResponse{
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
			RunresId:      item.RunresId,
			UsageDatetime: uint64(item.UsageDatetime.Unix()),
			Fee:           item.Fee,
			Discount:      item.Discount,
			InstanceId:    item.InstanceId,
			Type:          item.Type,
			InstanceName:  item.InstanceName,
			ResourceName:  item.ResourceName,
			Daynum:        item.Daynum,
			Minnum:        item.Minnum,
		}
	}

	return &types.MinuteUsageListResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: types.MinuteUsageListResponseData{
			Items: respItems,
			Total: total,
		},
	}, nil
}
