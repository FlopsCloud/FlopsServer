package logic

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/oauth2"
)

type XLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewXLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XLoginLogic {
	return &XLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XLoginLogic) XLogin(req *types.XLoginRequest) (resp *types.XLoginResponse, err error) {
	clientID := " "
	clientSecret := " -S"

	if req.RedirectURL == "" {
		req.RedirectURL = "https://hub.flopscloud.ai/api/v1/x/callback"
	}

	// 配置OAuth2客户端
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  req.RedirectURL, // 回调URL
		Scopes: []string{
			"tweet.read", "tweet.write", "users.read", "offline.access",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://twitter.com/i/oauth2/authorize",
			TokenURL: "https://api.twitter.com/2/oauth2/token",
		},
	}

	// state := "state" // In production, use a proper random state and store it in session/cookie

	// url := config.AuthCodeURL(state)
	url := buildAuthorizationURL(config)

	return &types.XLoginResponse{AuthURL: url}, nil

	// url := buildAuthorizationURL(config)
}

var codeVerifier string

func buildAuthorizationURL(config oauth2.Config) string {

	// PKCE  https://datatracker.ietf.org/doc/html/rfc7636

	codeVerifier = generateBase64Encoded32byteRandomString()
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	hashed := h.Sum(nil)
	codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashed)

	url := config.AuthCodeURL(
		"state",
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	return url
}
