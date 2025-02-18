package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/jwtx"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type WechatCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatCallbackLogic {
	return &WechatCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}

}

func (l *WechatCallbackLogic) FindOrCreateUser(svcCtx *svc.ServiceContext, openid string, nickname string) (*model.Users, error) {
	ctx := context.Background()

	user, err := svcCtx.UserModel.FindByOpenId(ctx, openid)

	if err == sqlx.ErrNotFound {

		userId, err := l.svcCtx.UserModel.RegisterNewUserinDB(l.ctx,
			nickname,
			nickname,
			openid+"@weixin.qq.com",
			"",
			"",
			"",
			openid,
			"",
			"",
			l.svcCtx.Config.Salt,
		)

		if err != nil {
			return nil, err
		}

		user, err = l.svcCtx.UserModel.FindOne(ctx, userId)
		if err != nil {
			return nil, err
		}

		return user, nil

	} else if err != nil {
		logx.Error("Failed to find user", err)
		return nil, err
	}
	return user, nil

}

func (l *WechatCallbackLogic) WechatCallback(req *types.WechatCallbackRequest) (resp *types.WechatCallbackResponse, err error) {
	code := req.Code

	// 用code、AppId和AppSecret向微信服务器换取用户信息
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", l.svcCtx.Config.WeixinOpen.AppId, l.svcCtx.Config.WeixinOpen.AppSecret, code)
	access_token_resp, err := http.Get(url)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	defer access_token_resp.Body.Close()
	body, err := ioutil.ReadAll(access_token_resp.Body)
	if err != nil {
		l.Logger.Error(err)

		return nil, err
	}
	l.Logger.Info(string(body))

	// {\"access_token\":\"86_i_NeH_t07_nUxSXWe2kHJ43WHF-YH8ym9DZ2YFq07ChS6IH6LLdyYL5Cjy48k7qgOKi5s3_eymWiYOKX0hVusLxumO4o9qfxN-I1XvqGoOQ\"
	// ,\"expires_in\":7200,
	// \"refresh_token\":\"86_M8jZNT_mOxAwpQ_jvO5PVjhrJ8eGzw6xW348YUU0IHxhDJQAOJRJo4YKMHngH-XNOP6OTycjlbHY5VnR-KrlDamObtqIXNxaT0vRop2JchA\",
	// \"openid\":\"oHzY96hsRBGClAVknigOSMWfSbvw\",
	// \"scope\":\"snsapi_login\",
	// \"unionid\":\"o4YUkwPVrBxtr4MU8MKW4FPlGb3U\"}"

	var accessToken AccessTokenResponse
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		l.Logger.Error(err)

		return nil, err
	}

	userInfo, err := getUserInfo(accessToken.OpenID, accessToken.AccessToken)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	//查找或创建用户

	user, err := l.FindOrCreateUser(l.svcCtx, userInfo.OpenID, userInfo.Nickname)
	if err != nil {
		l.Logger.Error(err)
		return &types.WechatCallbackResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Fail",
				Info:    err.Error(),
			}}, nil
	}

	// Generate JWT token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	token, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, (user.UserId), user.Email, user.SysRole)
	if err != nil {
		return &types.WechatCallbackResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "GetToken Fail",
				Info:    err.Error(),
			}}, nil
	}

	// openId, ok :=  result["openid"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("无法获取openid")
	// }
	// nickname, ok := result["nickname"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("无法获取nickname")
	// }
	// avatarUrl, ok := result["avatar_url"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("无法获取avatar_url")
	// }
	// 返回用户信息
	return &types.WechatCallbackResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		OpenId:    userInfo.OpenID,
		Nickname:  userInfo.Nickname,
		AvatarUrl: userInfo.HeadImgURL,
		Data: types.LoginResponseData{
			AccessToken:  token,
			AccessExpire: accessExpire,
		},
	}, nil

}

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

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UnionID      string `json:"unionid"`
	ExpiresIn    int    `json:"expires_in"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

func getAccessToken(AppId string, AppSecret string, code string) (*AccessTokenResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", AppId, AppSecret, code)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var accessTokenResponse AccessTokenResponse
	if err := json.Unmarshal(body, &accessTokenResponse); err != nil {
		return nil, err
	}

	return &accessTokenResponse, nil
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

	var userInfo WeChatUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
