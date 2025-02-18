package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListResourceUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListResourceUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListResourceUserLogic {
	return &ListResourceUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListResourceUserLogic) ListResourceUser(req *types.ListResourceUserRequest) (resp *types.ListResourceUserResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
