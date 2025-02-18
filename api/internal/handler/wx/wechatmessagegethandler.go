package wx

import (
	"net/http"

	"fca/api/internal/svc"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/zeromicro/go-zero/core/logc"
)

func WechatMessageGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wc := wechat.NewWechat()
		memory := cache.NewMemory()
		cfg := &config.Config{
			AppID:          svcCtx.Config.WeixinFWH.AppId,
			AppSecret:      svcCtx.Config.WeixinFWH.AppSecret,
			Token:          svcCtx.Config.WeixinFWH.Token,
			EncodingAESKey: svcCtx.Config.WeixinFWH.EncodingAESKey,
			Cache:          memory,
		}
		officialAccount := wc.GetOfficialAccount(cfg)

		// 传入request和responseWriter
		server := officialAccount.GetServer(r, w)

		logc.Info(r.Context(), r.Body)

		// 设置接收消息的处理方法
		server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {

			// 回复消息：演示回复用户发送的消息
			text := message.NewText(msg.Content)
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		})

		// 处理消息接收以及回复
		err := server.Serve()
		if err != nil {
			logc.Error(r.Context(), "微信消息处理失败", "error", err)

			return
		}
		// 发送回复的消息
		server.Send()

		// l := wx.NewWechatMessageGetLogic(r.Context(), svcCtx)
		// resp, err := l.WechatMessageGet()
		// if err != nil {
		// 	httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call" , err.Error()))
		// } else {
		// 	httpx.OkJsonCtx(r.Context(), w, resp)
		// }
	}
}
