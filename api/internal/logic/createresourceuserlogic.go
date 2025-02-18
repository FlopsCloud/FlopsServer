package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateResourceUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateResourceUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateResourceUserLogic {
	return &CreateResourceUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateResourceUserLogic) CreateResourceUser(req *types.CreateResourceUserRequest) (resp *types.ResourceUserResponse, err error) {
	// todo: add your logic here and delete this line

	// TODO: 根据用户请求注册用户

	return
}
