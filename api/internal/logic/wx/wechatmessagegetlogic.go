package wx

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WechatMessageGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatMessageGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatMessageGetLogic {
	return &WechatMessageGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatMessageGetLogic) WechatMessageGet() (resp *types.WechatMessageResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
