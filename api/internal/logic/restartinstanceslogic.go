package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

func RspNew(code int64, message string, info string) *types.Response {
	return &types.Response{
		Code:    code,
		Message: message,
		Info:    info,
	}

}

type RestartInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 所有 pod 能重新启动
func NewRestartInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestartInstancesLogic {
	return &RestartInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestartInstancesLogic) RestartInstances(req *types.RestartInstancesRequest, jwtToken string) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	// 2. is uid's

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

	// 检查余额
	balance, err := l.svcCtx.BalancesModel.FindOneByUserAndCurrency(l.ctx, userId, instance.OrgId, "CNY")
	if err != nil {
		resp.Code = response.NotFoundCode
		resp.Info = err.Error()
		resp.Message = "找不到余额"
		return resp, nil
	}
	if balance.Balance <= 0 {
		resp.Code = response.ServerErrorCode
		resp.Message = "余额不足,启动失败"
		resp.Info = fmt.Sprintf("余额: %f", float64(balance.Balance)/100)
		return resp, nil
	}

	// 检查实例是否到期
	expireDate := instance.ExpireDate
	if expireDate.Before(time.Now()) {
		resp.Code = response.ServerErrorCode
		resp.Message = "实例已到期,启动失败"
		return resp, nil
	}

	// 重启实例

	var message string
	if instance.IsChanged == 0 {
		resp2, err := RestartVhost(l.ctx, jwtToken, &DelVhostRequest{
			Name: instance.Name,
		})

		if err != nil {
			return RspNew(response.ServerErrorCode, "重启vhost失败", err.Error()), nil
		}

		if resp2.Code != 0 {
			return RspNew(response.ServerErrorCode, resp2.Message, resp2.Info), nil
		}
		message = resp2.Message
		l.svcCtx.InstancesLogModel.Insert(l.ctx, &model.InstancesLog{
			InstanceId: instance.InstanceId,
			UserId:     userId,
			Action:     "用户重启实例",
		})
	} else {

		DeleteVhost(l.ctx, jwtToken, &DelVhostRequest{
			Name: instance.Name,
		})

		image, err := l.svcCtx.ImagesModel.FindOne(l.ctx, instance.ImageId)
		if err != nil {
			return RspNew(response.ServerErrorCode, "找不到实例", err.Error()), nil
		}
		imageName := image.ImageName

		var resvhost *Response

		time.Sleep(2 * time.Second)
		for i := 0; i < 2; i++ {

			resvhost, err = CreateVhost(l.ctx, jwtToken, &VhostRequest{
				FcbPod: FcbPod{
					Name:    instance.Name,
					Port:    int32(instance.Port),
					Cpu:     instance.CpuCores,
					Memory:  instance.Memory,
					Storage: instance.Storage,
					Gpu:     instance.GpuCores,
					Image:   imageName,
					Mount:   instance.MountPath,
				},
			})
			if resvhost.Code == 500 {
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
		if err != nil {
			return RspNew(response.ServerErrorCode, "创建vhost失败", err.Error()), nil
		}

		if resvhost.Code != 0 {
			return RspNew(response.ServerErrorCode, resvhost.Message, resvhost.Info), nil
		}

		instance.IsChanged = 0
		l.svcCtx.InstanceModel.Update(l.ctx, instance)
		message = resvhost.Message

		l.svcCtx.InstancesLogModel.Insert(l.ctx, &model.InstancesLog{
			InstanceId: instance.InstanceId,
			UserId:     userId,
			Action:     "用户变更实例后重启",
		})

	}

	// 1. check instance status, is expire ?
	// 3. del pod
	// 4. waite pod delete
	// 5. new a pod

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "success",
		Info:    message,
	}, nil
}
