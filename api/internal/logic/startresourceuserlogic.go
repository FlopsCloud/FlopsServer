package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartResourceUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartResourceUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartResourceUserLogic {
	return &StartResourceUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartResourceUserLogic) StartResourceUser(req *types.StartResourceUserRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
