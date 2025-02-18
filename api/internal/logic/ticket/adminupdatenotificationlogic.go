package ticket

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateNotificationLogic {
	return &AdminUpdateNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateNotificationLogic) AdminUpdateNotification(req *types.AdminUpdateNotificationRequest) (resp *types.Response, err error) {
	resp = &types.Response{}

	role, _ := l.ctx.Value("role").(string)
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	if role != "superadmin" {
		resp.Code = response.UnauthorizedCode
		resp.Message = "只有超级管理员可以更新通知"
		resp.Info = "你的角色是" + role + ",id " + strconv.FormatInt(uid, 10)
		return resp, nil
	}

	notification, err := l.svcCtx.SystemNotificationsModel.FindOne(l.ctx, req.NotificationId)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Message = "通知不存在"
			return resp, nil
		}
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}

	notification.Title = req.Title
	notification.Content = req.Content
	notification.Type = req.Type
	if req.UserId > 0 {
		notification.UserId = sql.NullInt64{Int64: int64(req.UserId), Valid: true}
	}

	if req.OrgId > 0 {
		notification.OrgId = sql.NullInt64{Int64: int64(req.OrgId), Valid: true}
	}

	err = l.svcCtx.SystemNotificationsModel.Update(l.ctx, notification)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "更新通知失败"
		resp.Info = err.Error()
		return resp, nil
	}

	resp.Code = response.SuccessCode
	resp.Message = "更新通知成功"
	return resp, nil
}
