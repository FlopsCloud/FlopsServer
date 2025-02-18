package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDailyUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDailyUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDailyUsageLogic {
	return &ListDailyUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDailyUsageLogic) ListDailyUsage(req *types.DailyUsageListRequest) (resp *types.DailyUsageListResponse, err error) {
	// Ensure valid pagination
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Prepare optional filters
	var orgId, userId *uint64
	var startDate, endDate *string

	if req.OrgId != 0 {
		orgId = &req.OrgId
	}
	if req.UserId != 0 {
		userId = &req.UserId
	}
	if req.StartDate != "" {
		startDate = &req.StartDate
	}
	if req.EndDate != "" {
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
			UsageDate:       item.UsageDate.Format("2006-01-02"),
			ResourceId:      item.ResourceId,
			UsageMinAmount:  item.UsageMinAmount,
			UsageHourAmount: item.UsageHourAmount,
			UnitHourPrice:   item.UnitHourPrice,
			DiscountId:      item.DiscountId,
			Discount:        item.Discount,
		}
	}

	return &types.DailyUsageListResponse{
		Response: types.Response{
			Code:    0,
			Message: "success",
		},
		Data: types.DailyUsageListResponseData{
			Items: respItems,
			Total: total,
		},
	}, nil
}
