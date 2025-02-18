package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopResourceUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopResourceUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopResourceUserLogic {
	return &StopResourceUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopResourceUserLogic) StopResourceUser(req *types.StopResourceUserRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
