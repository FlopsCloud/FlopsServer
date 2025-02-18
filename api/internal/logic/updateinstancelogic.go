package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInstanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInstanceLogic {
	return &UpdateInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInstanceLogic) UpdateInstance(req *types.UpdateInstanceRequest) (resp *types.CreateInstanceResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
