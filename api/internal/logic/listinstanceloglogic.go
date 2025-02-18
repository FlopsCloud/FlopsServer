package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInstanceLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInstanceLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInstanceLogLogic {
	return &ListInstanceLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInstanceLogLogic) ListInstanceLog(req *types.ListInstanceLogRequest) (resp *types.ListInstanceLogResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
