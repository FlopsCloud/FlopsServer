package logic

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/jwtx"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/oauth2"
)

type XCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewXCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XCallbackLogic {
	return &XCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XCallbackLogic) XCallback(req *types.XCallbackRequest) (resp *types.XCallbackResponse, err error) {
	clientID := " "
	clientSecret := " - "
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://hub.flopscloud.ai/api/v1/x/callback", // 回调URL
		Scopes: []string{
			"tweet.read", "tweet.write", "users.read", "offline.access",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://twitter.com/i/oauth2/authorize",
			TokenURL: "https://api.twitter.com/2/oauth2/token",
		},
	}

	// token, err := config.Exchange(context.Background(), queryCode, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	// if err != nil {
	// 	log.Printf("failed to exchange token: %v\n", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// log.Printf("token scope: %v\n", token.Extra("scope"))

	// w.Header().Set("Content-Type", "application/json")
	// _, _ = io.Copy(w, res.Body)

	// Exchange the authorization code for tokens
	token, err := config.Exchange(l.ctx, req.Code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	// Get user info using the access token
	oAuthClient := oauth2.NewClient(l.ctx, oauth2.StaticTokenSource(token))

	// https://developer.x.com/en/docs/x-api/users/lookup/api-reference/get-users-me

	res, err := oAuthClient.Get("https://api.twitter.com/2/users/me?user.fields=profile_image_url")
	if err != nil {
		return nil, fmt.Errorf("failed to get me: %v", err)

	}
	defer res.Body.Close()

	userData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse user info
	var userInfo struct {
		Data struct {
			Id              string `json:"id"`
			Username        string `json:"username"`
			Name            string `json:"name"`
			ProfileImageUrl string `json:"profile_image_url"`
		} `json:"data"`
	}
	logx.Infof("userData: %s", string(userData))

	if err := json.Unmarshal(userData, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	user, err := l.svcCtx.UserModel.FindOneByX(l.ctx, userInfo.Data.Id)
	if err != nil {
		// 创建用户

		//TODO email 可能重复，这个时候应该用随机数

		user = &model.Users{
			Username:     userInfo.Data.Name,
			Nickname:     userInfo.Data.Username,
			Email:        userInfo.Data.Name + "@x.com",
			Phone:        "",
			PasswordHash: "",
			AccX:         userInfo.Data.Id,
			HeadUrl:      userInfo.Data.ProfileImageUrl,
		}
		res, err := l.svcCtx.UserModel.Insert(l.ctx, user)
		if err != nil {
			l.Logger.Error("XCallback create user error", err)
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
	return &types.XCallbackResponse{
		AccessToken: jwtToken,
		// RefreshToken: token.RefreshToken,
		AccessExpire: now + accessExpire,
		Raw:          string(userData),
		UserInfo: &types.XUserInfo{
			Id:            userInfo.Data.Id,
			Username:      userInfo.Data.Username,
			Name:          userInfo.Data.Name,
			Picture:       userInfo.Data.ProfileImageUrl,
			VerifiedEmail: false,
		},
	}, nil
}
func generateBase64Encoded32byteRandomString() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}
