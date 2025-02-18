package ticket

import (
	"context"
	"database/sql"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotificationLogic {
	return &CreateNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNotificationLogic) CreateNotification(req *types.CreateNotificationRequest) (resp *types.NotificationResponse, err error) {

	resp = &types.NotificationResponse{}

	role, _ := l.ctx.Value("role").(string)

	if role != "superadmin" && role != "admin" {
		resp.Code = response.UnauthorizedCode
		resp.Message = "只有管理员可以创建通知"
		return resp, nil
	}
	// 管理员权限验证 end

	notification := &model.SystemNotifications{
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
	}

	// Handle optional fields
	if req.UserId != 0 {
		notification.UserId = sql.NullInt64{
			Int64: int64(req.UserId),
			Valid: true,
		}
	}
	if req.OrgId != 0 {
		notification.OrgId = sql.NullInt64{
			Int64: int64(req.OrgId),
			Valid: true,
		}
	}

	result, err := l.svcCtx.SystemNotificationsModel.Insert(l.ctx, notification)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}

	// Return success response with created notification
	resp.Code = response.SuccessCode
	resp.Message = "Success"
	resp.Data = types.SystemNotification{
		NotificationId: uint64(id),
		Title:          notification.Title,
		Content:        notification.Content,
		Type:           notification.Type,
		// Status:         notification.Status,
		CreatedAt: uint64(time.Now().Unix()),
	}

	if notification.UserId.Valid {
		resp.Data.UserId = uint64(notification.UserId.Int64)
	}
	if notification.OrgId.Valid {
		resp.Data.OrgId = uint64(notification.OrgId.Int64)
	}

	return resp, nil
}
