package logic

import (
	"context"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WechatLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatLoginLogic {
	return &WechatLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatLoginLogic) WechatLogin(req *types.WechatLoginRequest) (resp *types.WechatLoginResponse, err error) {

	appId := l.svcCtx.Config.WeixinOpen.AppId
	redirectUri := req.RedirectUri
	scope := "snsapi_login"
	state := "STATE_STRING"

	qrConnectUrl := fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect", appId, redirectUri, scope, state)
	// 返回一个包含微信扫码授权链接的响应，让客户端去引导用户扫码
	return &types.WechatLoginResponse{
		Url:       qrConnectUrl,
		OpenId:    "",
		Nickname:  "尚未获取",
		AvatarUrl: "",
	}, nil
}
