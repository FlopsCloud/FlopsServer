package ticket

import (
	"context"
	"encoding/json"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkNotificationReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkNotificationReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkNotificationReadLogic {
	return &MarkNotificationReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkNotificationReadLogic) MarkNotificationRead(req *types.MarkNotificationReadRequest) (resp *types.Response, err error) {
	resp = &types.Response{}
	nid := req.NotificationId

	// Get current user ID from context
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		resp.Code = response.UnauthorizedCode
		resp.Message = err.Error()
		return resp, nil
	}
	userId := uint64(uid)

	// Convert notification ID from string to uint64

	// Check if notification exists TODO remove reading
	notification, err := l.svcCtx.SystemNotificationsModel.FindOne(l.ctx, nid)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Message = "Notification not found"
			return resp, nil
		}
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}
	logx.Info(notification)

	// Create notification read record
	readRecord := &model.NotificationReads{
		NotificationId: nid,
		UserId:         userId,
		IsRead:         1,
		ReadAt:         time.Now(),
	}

	_, err = l.svcCtx.NotificationReadsModel.Insert(l.ctx, readRecord)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}

	// Update notification status to "read"

	// notification.Status = "read"
	// err = l.svcCtx.SystemNotificationsModel.Update(l.ctx, notification)
	// if err != nil {
	// 	resp.Code = response.ServerErrorCode
	// 	resp.Message = err.Error()
	// 	return resp, nil
	// }

	resp.Code = response.SuccessCode
	resp.Message = "Success"
	return resp, nil
}
