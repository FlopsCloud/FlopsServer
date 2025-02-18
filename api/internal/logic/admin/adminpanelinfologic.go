package admin

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminPanelInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminPanelInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminPanelInfoLogic {
	return &AdminPanelInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminPanelInfoLogic) AdminPanelInfo() (resp *types.AdminPanelInfoResp, err error) {
	// Check admin permissions
	role, _ := l.ctx.Value("role").(string)
	if role != "admin" && role != "superadmin" {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.UnauthorizedCode,
				Message: "only admin can access",
			},
		}, nil
	}

	// Get total users
	totalUsers, err := l.svcCtx.UserModel.Count(l.ctx, "is_deleted = 0")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get users count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get total instances
	totalInstances, err := l.svcCtx.InstanceModel.Count(l.ctx, "1=1")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get instances count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get running instances
	runningInstances, err := l.svcCtx.InstanceModel.Count(l.ctx, "state = 'running'")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get running instances count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get total organizations
	totalOrgs, err := l.svcCtx.OrganizationModel.Count(l.ctx, "  1 = 1")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get organizations count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get total servers
	totalServers, err := l.svcCtx.ServerModel.Count(l.ctx, "where 1=1")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get servers count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get total orders and unpaid orders
	totalOrders, err := l.svcCtx.OrderRecordsModel.Count(l.ctx, "where 1=1")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get orders count",
				Info:    err.Error(),
			},
		}, nil
	}

	unpaidOrders, err := l.svcCtx.OrderRecordsModel.Count(l.ctx, "where 1=1")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get unpaid orders count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get recent instances (limit 5)
	recentInstances, err := l.svcCtx.InstanceModel.FindAll(l.ctx, "1=1 order by created_at desc limit 5")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get recent instances",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get recent orders (limit 5)
	recentOrders, err := l.svcCtx.OrderRecordsModel.FindList(l.ctx, "where 1=1 order by created_at desc limit 5")
	if err != nil {
		return &types.AdminPanelInfoResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get recent orders",
				Info:    err.Error(),
			},
		}, nil
	}

	// Convert instances to response type
	var instanceList []types.Instance
	for _, instance := range recentInstances {
		instanceList = append(instanceList, types.Instance{
			InstanceId: instance.InstanceId,
			Name:       instance.Name,
			State:      instance.State,
			CreatedAt:  uint64(instance.CreatedAt.Unix()),
		})
	}

	// Convert orders to response type
	var orderList []types.OrderRecords
	for _, order := range *recentOrders {
		orderList = append(orderList, types.OrderRecords{
			OrderNo:     order.OrderNo,
			UserId:      order.UserId,
			OrderAmount: order.OrderAmount,
			CreatedAt:   uint64(order.CreatedAt.Unix()),
		})
	}

	return &types.AdminPanelInfoResp{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: types.AdminPanelInfoRespData{
			TotalUser:        totalUsers,
			TotalOrg:         totalOrgs,
			TotalServer:      uint64(totalServers),
			TotalInstance:    totalInstances,
			RunningInstances: runningInstances,
			TotalOrder:       uint64(totalOrders),
			UnpaidOrder:      uint64(unpaidOrders),
			UnreadMessage:    0, // These could be implemented later if needed
			UnreadTicket:     0,
			RecentInstance:   instanceList,
			RecentOrder:      orderList,
			// DailyUsage:       []types.DailyUsageItem{}, // This could be implemented later
		},
	}, nil
}
