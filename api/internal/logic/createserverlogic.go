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

type CreateServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateServerLogic {
	return &CreateServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateServerLogic) CreateServer(req *types.CreateServerRequest) response.Response {
	expireDate, err := time.Parse("2006-01-02", req.ExpireDate)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	res, err := l.svcCtx.ServerModel.Insert(l.ctx, &model.Servers{
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
		PayPrices:     "{}",
		DiskPayPrices: "{}",
		IsOnline:      1,
		IsPayDaily:    1,
		IsPayMin:      1,
		IsPayMonthly:  1,
		IsPayYearly:   1,
		CpuFrequency:  2500,
	})
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	// req.ServerId = uint64(id)
	return response.OK(&types.Server{
		ServerId:       uint64(id),
		BillingMethod:  req.BillingMethod,
		Region:         req.Region,
		Supplier:       req.Supplier,
		Processor:      req.Processor,
		GpuModel:       req.GpuModel,
		GpuCount:       req.GpuCount,
		GraphicsMemory: req.GraphicsMemory,
		CpuModel:       req.CpuModel,
		CpuCores:       req.CpuCores,
		ProcessorCount: req.ProcessorCount,
		Memory:         req.Memory,
		SystemDisk:     req.SystemDisk,
		DataDisk:       req.DataDisk,
		MaxDataDisk:    req.MaxDataDisk,
		SpeedDesc:      req.SpeedDesc,
		GpuDriver:      req.GpuDriver,
		CudaVersion:    req.CudaVersion,
		Cost:           req.Cost,
		Available:      req.Available,
		Total:          req.Total,
		ExpireDate:     req.ExpireDate,
	})
}
