package wx

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"time"

	"github.com/google/uuid"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

// 微信登录配置
type WechatConfig struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
}

type GenerateQRCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateQRCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateQRCodeLogic {
	return &GenerateQRCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateQRCodeLogic) GenerateQRCode(req *types.GenerateQRCodeReq) (resp *types.GenerateQRCodeResp, err error) {
	// 1. 创建唯一票据ID
	// ticketId := generateTicketId()

	// 2. 初始化微信实例

	resp = &types.GenerateQRCodeResp{}

	wechatCfg := &config.Config{
		AppID:          l.svcCtx.Config.WeixinFWH.AppId,
		AppSecret:      l.svcCtx.Config.WeixinFWH.AppSecret,
		Token:          l.svcCtx.Config.WeixinFWH.Token,
		EncodingAESKey: l.svcCtx.Config.WeixinFWH.EncodingAESKey,
	}

	officialAccount := wechatClient.GetOfficialAccount(wechatCfg)

	Scene := uuid.New().String()

	basicObj := officialAccount.GetBasic()
	tq := basic.NewTmpQrRequest(time.Duration(60*5), Scene)
	ticket, err := basicObj.GetQRTicket(tq)
	if err != nil {
		logc.Errorf(l.ctx, "get qr ticket failed, %s", err)
		resp.Code = response.ServerErrorCode
		resp.Message = "get qr ticket failed"
		resp.Info = err.Error()

		return resp, nil
	}

	// 3. 保存到Redis中
	l.svcCtx.RedisClient.Setex("QR_"+Scene, "WAITING_SCAN", int(ticket.ExpireSeconds*2))
	logc.Info(l.ctx, "ticket: ", ticket.Ticket, " expire:", ticket.ExpireSeconds)

	url := basic.ShowQRCode(ticket)

	resp.Code = response.SuccessCode
	resp.Message = "success"

	resp.Data = types.QRCodeData{
		TicketId:  ticket.Ticket,
		QrcodeUrl: url,
		ExpireAt:  time.Now().Add(time.Second * time.Duration(ticket.ExpireSeconds)).Unix(),
		Scene:     Scene,
	}

	return resp, nil

}
