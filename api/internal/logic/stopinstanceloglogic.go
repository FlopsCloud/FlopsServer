package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopInstanceLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopInstanceLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopInstanceLogLogic {
	return &StopInstanceLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopInstanceLogLogic) StopInstanceLog(req *types.StopInstanceLogRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
