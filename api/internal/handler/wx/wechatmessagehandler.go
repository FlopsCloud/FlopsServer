package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"fca/api/internal/svc"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/zeromicro/go-zero/core/logc"
)

var (
	wc              *wechat.Wechat
	officialAccount *officialaccount.OfficialAccount
)

// func getUserInfo(fromUserName string) (*user.Info, error) {
// 	userInfoManager := officialAccount.GetUser()

// 	// Retrieve user information
// 	userInfo, err := userInfoManager.GetUserInfo(fromUserName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logc.Info(nil, "userinfo", userInfo)

// 	return userInfo, err
// }

type WeChatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	AppID        string `json:"appid"`
	Secret       string `json:"secret"`
	ForceRefresh bool   `json:"force_refresh,omitempty"` // 可选参数
}

func getStableAccessToken(appID string, appSecret string, forceRefresh bool) (*TokenResponse, error) {
	url := "https://api.weixin.qq.com/cgi-bin/stable_token"
	requestBody := TokenRequest{
		GrantType:    "client_credential",
		AppID:        appID,
		Secret:       appSecret,
		ForceRefresh: forceRefresh,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func getUserInfo(openID, accessToken string) (*WeChatUserInfo, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", accessToken, openID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logc.Info(nil, "userinfo", string(body))

	var userInfo WeChatUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func WechatMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if officialAccount == nil {
			wc = wechat.NewWechat()
			memory := cache.NewMemory()
			cfg := &config.Config{
				AppID:          svcCtx.Config.WeixinFWH.AppId,
				AppSecret:      svcCtx.Config.WeixinFWH.AppSecret,
				Token:          svcCtx.Config.WeixinFWH.Token,
				EncodingAESKey: svcCtx.Config.WeixinFWH.EncodingAESKey,
				Cache:          memory,
			}
			officialAccount = wc.GetOfficialAccount(cfg)

		}

		// wc := wechat.NewWechat()
		// memory := cache.NewMemory()
		// cfg := &config.Config{
		// 	AppID:          svcCtx.Config.WeixinFWH.AppId,
		// 	AppSecret:      svcCtx.Config.WeixinFWH.AppSecret,
		// 	Token:          svcCtx.Config.WeixinFWH.Token,
		// 	EncodingAESKey: svcCtx.Config.WeixinFWH.EncodingAESKey,
		// 	Cache:          memory,
		// }
		// officialAccount := wc.GetOfficialAccount(cfg)

		// 传入request和responseWriter
		server := officialAccount.GetServer(r, w)

		// logc.Info(r.Context(), r.Body)

		// 设置接收消息的处理方法
		server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {

			switch msg.MsgType {

			case message.MsgTypeEvent:
				logc.Info(r.Context(), "收到事件消息", msg)
				if msg.Event == "SCAN" {

					// token, err := getStableAccessToken(svcCtx.Config.WeixinFWH.AppId, svcCtx.Config.WeixinFWH.AppSecret, false)

					// // accesstoken, err := officialAccount.GetAccessToken()
					// if err != nil {
					// 	return nil
					// }

					// //msg.FromUserName
					// logc.Info(r.Context(), "accesstoken: ", token.AccessToken)
					// userInfo, err := getUserInfo(string(msg.FromUserName), token.AccessToken)
					// if err != nil {
					// 	logc.Info(r.Context(), "获取用户信息失败: %v", err)
					// 	return nil
					// }

					// text := message.NewText(userInfo.Nickname + " 欢迎登录")
					svcCtx.RedisClient.Setex("QR_"+msg.EventKey, string(msg.FromUserName), 300)
					text := message.NewText("欢迎登录")
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}

			case message.MsgTypeText:

				logc.Info(r.Context(), "收到文本消息", msg)
				text := message.NewText(msg.Content)
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}

			case message.MsgTypeImage:
				//处理图片消息
				//return handleImageMessage(msg)
			case message.MsgTypeVoice:
				// 处理语音消息
				//return handleVoiceMessage(msg)

			default:
				text := message.NewText("啥??")
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}

			}
			return nil

		})

		// 处理消息接收以及回复
		err := server.Serve()
		if err != nil {
			logc.Error(r.Context(), "微信消息处理失败", "error", err)

			return
		}
		// 发送回复的消息
		server.Send()

		// var req types.WechatMessageRequest
		// if err := httpx.Parse(r, &req); err != nil {
		// 	httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call" , err.Error()))
		// 	return
		// }

		// l := wx.NewWechatMessageLogic(r.Context(), svcCtx)
		// resp, err := l.WechatMessage(&req)
		// if err != nil {
		// 	httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call" , err.Error()))
		// } else {
		// 	httpx.OkJsonCtx(r.Context(), w, resp)
		// }
	}
}
