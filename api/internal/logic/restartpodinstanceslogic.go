package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RestartPodInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestartPodInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestartPodInstancesLogic {
	return &RestartPodInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestartPodInstancesLogic) RestartPodInstances(req *types.RestartInstancesRequest, jwtToken string) (resp *types.Response, err error) {
	resp = &types.Response{}

	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		resp.Code = response.UnauthorizedCode
		resp.Info = err.Error()
		resp.Message = "需要登录"
		return resp, nil
	}
	userId := uint64(uid)

	instance, err := l.svcCtx.InstanceModel.FindOne(l.ctx, req.InstanceId)
	if err != nil {
		resp.Code = response.NotFoundCode
		resp.Info = err.Error()
		resp.Message = "找不到实例"
		return resp, nil
	}
	if instance.UserId != userId {
		resp.Code = response.UnauthorizedCode
		resp.Message = "没有权限"
		return resp, nil
	}

	resp2, err := RestartVhost(l.ctx, jwtToken, &DelVhostRequest{
		Name: instance.Name,
	})

	if err != nil {
		return RspNew(response.ServerErrorCode, "重启vhost失败", err.Error()), nil
	}

	if resp2.Code != 0 {
		return RspNew(response.ServerErrorCode, resp2.Message, resp2.Info), nil
	}

	l.svcCtx.InstancesLogModel.Insert(l.ctx, &model.InstancesLog{
		InstanceId: instance.InstanceId,
		UserId:     userId,
		Action:     "用户重启POD实例",
	})
	instance.IsChanged = 0
	l.svcCtx.InstanceModel.Update(l.ctx, instance)

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "success",
		Info:    resp2.Message,
	}, nil
}
