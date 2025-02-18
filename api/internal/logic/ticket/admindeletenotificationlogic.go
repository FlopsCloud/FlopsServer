package ticket

import (
	"context"
	"encoding/json"
	"strconv"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteNotificationLogic {
	return &AdminDeleteNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteNotificationLogic) AdminDeleteNotification(req *types.AdminDeleteNotificationRequest) (resp *types.Response, err error) {
	resp = &types.Response{}

	role, _ := l.ctx.Value("role").(string)
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	if role != "superadmin" && role != "admin" {
		resp.Code = response.UnauthorizedCode
		resp.Message = "只有超级管理员和系统管理员可以删除通知"
		resp.Info = "你的角色是" + role + ",id " + strconv.FormatInt(uid, 10)
		return resp, nil
	}
	err = l.svcCtx.SystemNotificationsModel.Delete(l.ctx, req.NotificationId)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = response.NotFoundCode
			resp.Message = "通知不存在"
			return resp, nil
		}
		resp.Code = response.ServerErrorCode
		resp.Info = err.Error()
		resp.Message = "删除通知失败"
		return resp, nil
	}

	resp.Code = response.SuccessCode
	resp.Message = "删除通知成功"
	return resp, nil
}
