package ticket

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSystemMetricsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSystemMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSystemMetricsLogic {
	return &GetSystemMetricsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSystemMetricsLogic) GetSystemMetrics(req *types.GetSystemMetricsRequest) (resp *types.GetSystemMetricsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
