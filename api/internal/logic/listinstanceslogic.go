package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInstancesLogic {
	return &ListInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInstancesLogic) ListInstances(req *types.InstanceListRequest) (resp *types.InstanceListResponse) {
	resp = &types.InstanceListResponse{
		Response: types.Response{},
	}

	// Get instances with filters and pagination
	instances, total, err := l.svcCtx.InstanceModel.FindByFilter(l.ctx, req.UserId, req.ServerId, req.Page, req.PageSize)
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}

	// Convert model instances to response instances
	var instanceList []types.Instance
	for _, instance := range instances {
		instanceList = append(instanceList, types.Instance{
			InstanceId:    instance.InstanceId,
			Name:          instance.Name,
			State:         instance.State,
			UserId:        instance.UserId,
			ServerId:      instance.ServerId,
			ImageId:       instance.ImageId,
			Ip:            instance.Ip,
			Port:          instance.Port,
			ContainerPort: instance.ContainerPort,
			GpuCores:      instance.GpuCores,
			Memory:        instance.Memory,
			DiskPath:      instance.DiskPath,
			Cost:          instance.Cost,
			MountPath:     instance.MountPath,
		})
	}

	resp.Response = types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}
	resp.Data = types.InstanceListResponseData{
		Instances: instanceList,
		Total:     total,
	}

	return resp
}
