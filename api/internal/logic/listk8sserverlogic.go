package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListK8sServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListK8sServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListK8sServerLogic {
	return &ListK8sServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListK8sServerLogic) ListK8sServer(jwtToken string) (resp *types.ListServersResponse, err error) {

	var serverinfos []types.ServerInfo
	resp = &types.ListServersResponse{
		Servers: serverinfos,
	}
	servers, err := ListServers(l.ctx, jwtToken)
	if err != nil {
		return nil, err
	}

	// Name        string `json:"name"`
	// TotalCPU    string `json:"total_cpu"`
	// UsedCPU     string `json:"used_cpu"`
	// TotalMemory string `json:"total_memory"`
	// UsedMemory  string `json:"used_memory"`
	// TotalGPU    string `json:"total_gpu"`
	// UsedGPU     string `json:"used_gpu"`
	// GPUType     string `json:"gpu_type"`
	// Status      string `json:"status"`
	for _, server := range servers.Servers {
		serverinfos = append(serverinfos, types.ServerInfo{

			Name:        server.Name,
			TotalCPU:    server.TotalCPU,
			UsedCPU:     server.UsedCPU,
			TotalMemory: server.TotalMemory,
			UsedMemory:  server.UsedMemory,
			TotalGPU:    server.TotalGPU,
			UsedGPU:     server.UsedGPU,
			GPUType:     server.GPUType,
			Status:      server.Status,
		})
	}

	resp.Code = response.SuccessCode
	resp.Message = "success"
	resp.Servers = serverinfos
	return resp, nil
}
