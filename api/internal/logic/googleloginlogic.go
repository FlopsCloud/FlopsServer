package logic

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoogleLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGoogleLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoogleLoginLogic {
	return &GoogleLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// https://hub.flopscloud.ai/api/v1/oauth2/login

func (l *GoogleLoginLogic) GoogleLogin(req *types.GLoginRequest) (*types.GLoginResponse, error) {

	if req.RedirectURL == "" {
		req.RedirectURL = "https://hub.flopscloud.ai/api/v1/google/callback"
	}

	config := &oauth2.Config{
		ClientID:     l.svcCtx.Config.Google.Client,
		ClientSecret: l.svcCtx.Config.Google.Key,
		RedirectURL:  req.RedirectURL, // Update this with your actual callback URL
		Scopes: []string{
			"email", "profile", "openid",
			//"https://www.googleapis.com/auth/userinfo.email",
			//"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Generate random state
	//state := "random-state" // In production, use a proper random state and store it in session/cookie
	state := "state" // In production, use a proper random state and store it in session/cookie

	// Get the URL to redirect to Google's consent page
	url := config.AuthCodeURL(state)

	return &types.GLoginResponse{AuthURL: url}, nil
}
