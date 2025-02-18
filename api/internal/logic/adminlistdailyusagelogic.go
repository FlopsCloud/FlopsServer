package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListDailyUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListDailyUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListDailyUsageLogic {
	return &AdminListDailyUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListDailyUsageLogic) AdminListDailyUsage(req *types.AdminDailyUsageListRequest) (resp *types.DailyUsageListResponse, err error) {
	role, _ := l.ctx.Value("role").(string)

	if role != "superadmin" || role == "admin" {
		return &types.DailyUsageListResponse{
			Response: types.Response{
				Code:    response.UnauthorizedCode,
				Message: "permission denied",
			},
		}, nil
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Prepare optional filters
	var orgId, userId *uint64
	var startDate, endDate *uint64

	if req.OrgId != 0 {
		orgId = &req.OrgId
	}

	if req.UserId != 0 {
		userId = &req.UserId
	}

	if req.StartDate != 0 {
		startDate = &req.StartDate
	}
	if req.EndDate != 0 {
		endDate = &req.EndDate
	}

	// Get data from database
	items, total, err := l.svcCtx.DailyUsageModel.List(l.ctx, orgId, userId, startDate, endDate, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// Convert model items to response items
	respItems := make([]types.DailyUsageItem, len(items))
	for i, item := range items {
		respItems[i] = types.DailyUsageItem{
			UsageId:         item.UsageId,
			OrgId:           item.OrgId,
			UserId:          item.UserId,
			UsageDate:       uint64(item.UsageDate.Unix()),
			RunresId:        item.RunresId,
			UsageMinAmount:  item.UsageMinAmount,
			UsageHourAmount: item.UsageHourAmount,
			UnitHourPrice:   item.UnitHourPrice,
			DiscountId:      item.DiscountId,
			Discount:        item.Discount,
			InstanceId:      item.InstanceId,
			Type:            item.Type,
			Fee:             item.Fee,
			InstanceName:    item.InstanceName,
			ResourceName:    item.ResourceName,
			Daynum:          item.Daynum,
		}
	}

	return &types.DailyUsageListResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: types.DailyUsageListResponseData{
			Items: respItems,
			Total: total,
		},
	}, nil
}
