package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TerminateInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTerminateInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TerminateInstancesLogic {
	return &TerminateInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TerminateInstancesLogic) TerminateInstances(req *types.TerminateInstancesRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
