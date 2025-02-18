package ticket

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNotificationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationsLogic {
	return &ListNotificationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNotificationsLogic) ListNotifications(req *types.ListNotificationsRequest) (resp *types.ListNotificationsResponse, err error) {
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

	// Set default page size if not provided
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	// Get notifications using model method
	notificationsEx, total, err := l.svcCtx.SystemNotificationsModel.FindByFilter2(
		l.ctx,
		userId,
		req.OrgId,
		req.Type,
		req.Status, //unread,read,all
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
	for _, n := range notificationsEx {
		notification := types.SystemNotification{
			NotificationId: n.NotificationId,
			Title:          n.Title,
			Content:        n.Content,
			Type:           n.Type,
			IsRead:         n.IsRead,
			CreatedAt:      uint64(n.CreatedAt.Unix()),
			// Status:         n.Status,
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
