package wx

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WechatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatMessageLogic {
	return &WechatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatMessageLogic) WechatMessage(req *types.WechatMessageRequest) (resp *types.WechatMessageResponse, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Info(req)

	return
}
