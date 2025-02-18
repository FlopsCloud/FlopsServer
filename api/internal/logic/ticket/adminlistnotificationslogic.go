package ticket

import (
	"context"
	"encoding/json"
	"strconv"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListNotificationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListNotificationsLogic {
	return &AdminListNotificationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListNotificationsLogic) AdminListNotifications(req *types.ListNotificationsRequest) (resp *types.ListNotificationsResponse, err error) {
	resp = &types.ListNotificationsResponse{}

	// Get current user ID from context
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		resp.Code = response.UnauthorizedCode
		resp.Info = err.Error()
		resp.Message = "需要登录"
		return resp, nil
	}
	userId := uint64(uid)

	role, _ := l.ctx.Value("role").(string)
	if role != "superadmin" && role != "admin" {
		resp.Code = response.UnauthorizedCode
		resp.Info = "你的角色是" + role + ",id " + strconv.FormatInt(uid, 10)
		resp.Message = "只有超级管理员和系统管理员可以查看通知"
		return resp, nil
	}

	// Set default page size if not provided
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	// Get notifications using model method
	notifications, total, err := l.svcCtx.SystemNotificationsModel.FindByFilter(
		l.ctx,
		userId,
		req.OrgId,
		// 0,
		req.Type,
		// req.Status,
		int64(req.Page),
		int64(pageSize),
	)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}

	// Convert to response type
	notificationList := make([]types.SystemNotification, 0)
	for _, n := range notifications {
		notification := types.SystemNotification{
			NotificationId: n.NotificationId,
			Title:          n.Title,
			Content:        n.Content,
			Type:           n.Type,
			// Status:         n.Status,
			CreatedAt: uint64(n.CreatedAt.Unix()),
		}
		if n.UserId.Valid {
			notification.UserId = uint64(n.UserId.Int64)
		}
		if n.OrgId.Valid {
			notification.OrgId = uint64(n.OrgId.Int64)
		}
		notificationList = append(notificationList, notification)
	}

	resp.Code = response.SuccessCode
	resp.Message = "Success"
	resp.Data = types.ListNotificationsResponseData{
		Notifications: notificationList,
		Total:         total,
	}

	return resp, nil

}
