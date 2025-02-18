package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopInstancesLogic {
	return &StopInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopInstancesLogic) StopInstances(req *types.StopInstancesRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
