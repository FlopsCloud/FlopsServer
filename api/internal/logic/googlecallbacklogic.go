package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/jwtx"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoogleCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGoogleCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoogleCallbackLogic {
	return &GoogleCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoogleCallbackLogic) GoogleCallback(req *types.CallbackRequest) (*types.CallbackResponse, error) {
	config := &oauth2.Config{
		ClientID:     l.svcCtx.Config.Google.Client,
		ClientSecret: l.svcCtx.Config.Google.Key,
		RedirectURL:  "https://hub.flopscloud.ai/api/v1/google/callback", // Update this with your actual callback URL
		Scopes: []string{
			"email", "profile", "openid",
			//"https://www.googleapis.com/auth/userinfo.email",
			//"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Exchange the authorization code for tokens
	token, err := config.Exchange(l.ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	// Get user info using the access token
	client := config.Client(l.ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	userData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse user info
	var userInfo struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	if err := json.Unmarshal(userData, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	user, err := l.svcCtx.UserModel.FindOneByGoogle(l.ctx, userInfo.Email)
	if err != nil {
		// 创建用户
		user = &model.Users{
			Username:     userInfo.Email,
			Nickname:     userInfo.Name,
			Email:        userInfo.Email,
			Phone:        "",
			PasswordHash: "",
			AccGoogle:    userInfo.Email,
			HeadUrl:      userInfo.Picture,
		}

		res, err := l.svcCtx.UserModel.Insert(l.ctx, user)
		if err != nil {
			l.Logger.Error("GoogleCallback create user error", err)
			return nil, err
		}
		userId, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		user.UserId = uint64(userId)
	}

	//login success

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, int64(user.UserId), user.Email)
	if err != nil {
		l.Logger.Error("GoogleCallback create user error", err)
		return nil, err
	}

	// Create response

	return &types.CallbackResponse{
		AccessToken: jwtToken,
		// RefreshToken: token.RefreshToken,
		AccessExpire: now + accessExpire,
		UserInfo: &types.UserInfo{
			Id:            userInfo.Id,
			Email:         userInfo.Email,
			Name:          userInfo.Name,
			Picture:       userInfo.Picture,
			VerifiedEmail: userInfo.VerifiedEmail,
		},
	}, nil
}
