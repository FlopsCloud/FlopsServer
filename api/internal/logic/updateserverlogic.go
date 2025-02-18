package logic

import (
	"context"
	"fca/common/response"
	"fca/model"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateServerLogic {
	return &UpdateServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateServerLogic) UpdateServer(req *types.Server) response.Response {
	expireDate, err := time.Parse("2006-01-02", req.ExpireDate)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	err = l.svcCtx.ServerModel.Update(l.ctx, &model.Servers{
		ServerId:      req.ServerId,
		BillingMethod: req.BillingMethod,
		RegionId:      req.Region,
		Supplier:      req.Supplier,
		Processor:     req.Processor,
		GpuModel:      req.GpuModel,
		GpuCount:      req.GpuCount,
		GpuMem:        req.GraphicsMemory,
		CpuModel:      req.CpuModel,
		CpuCores:      req.CpuCores,
		ThreadCount:   req.ProcessorCount,
		Memory:        req.Memory,
		SystemDisk:    req.SystemDisk,
		DataDisk:      req.DataDisk,
		MaxDataDisk:   req.MaxDataDisk,
		SpeedDesc:     req.SpeedDesc,
		GpuDriver:     req.GpuDriver,
		CudaVersion:   req.CudaVersion,
		Cost:          req.Cost,
		Available:     req.Available,
		Total:         req.Total,
		ExpireDate:    expireDate,
	})
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	return response.OK(req)
}
