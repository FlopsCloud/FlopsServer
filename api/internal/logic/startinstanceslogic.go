package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartInstancesLogic {
	return &StartInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartInstancesLogic) StartInstances(req *types.StartInstancesRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
