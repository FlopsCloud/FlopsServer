package wx

import (
	"context"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/jwtx"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CheckQRCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckQRCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckQRCodeLogic {
	return &CheckQRCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func LoginErr(code int64, msg string, err error) *types.LoginResponse {
	var info string
	if err != nil {
		info = err.Error()
	}
	return &types.LoginResponse{
		Response: types.Response{
			Code:    code,
			Message: msg,
			Info:    info,
		},
	}
}

func (l *CheckQRCodeLogic) CheckQRCode(req *types.CheckQRCodeReq) (resp *types.LoginResponse, err error) {
	// 检查key是否存在
	key := req.Scene
	if key == "" {
		return LoginErr(400, "Scene is empty", nil), nil
	}

	loop_cnt := 0
	var openid string

	for {
		loop_cnt++
		if loop_cnt > 30 {
			return LoginErr(300, "Time out, Key is not set yet, Retry", nil), nil
		}

		openid, err = l.svcCtx.RedisClient.Get("QR_" + key)
		if err != nil {
			return LoginErr(response.ServerErrorCode, "Get Redis Error", err), nil

		}

		if openid == "" {
			return LoginErr(400, "key is not found", err), nil
		}

		if openid != "WAITING_SCAN" {
			break
		}
		time.Sleep(time.Second)
	}

	logx.Info("key is found ", key, " ", openid)
	// Fetch user from database
	user, err := l.svcCtx.UserModel.FindByOpenId(l.ctx, openid)

	var userid uint64

	if err == sqlx.ErrNotFound {

		userId, err := l.svcCtx.UserModel.RegisterNewUserinDB(l.ctx,
			"微信用户",
			"微信用户",
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
			return &types.LoginResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Failed to Create user",
					Info:    err.Error(),
				},
			}, nil
		}

		user, err = l.svcCtx.UserModel.FindOne(l.ctx, userId)

		if err != nil {
			l.Logger.Info("CheckQRCode Failed to insert user", err)

			return &types.LoginResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Failed to Create user",
					Info:    err.Error(),
				},
			}, nil
		}

		userid = user.UserId

	} else if err != nil {
		return &types.LoginResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to fetch user",
				Info:    err.Error(),
			},
		}, nil
	} else {
		userid = user.UserId
	}

	//TODO , get userinfo from wechat

	// Generate JWT token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	token, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, userid, user.Email, user.SysRole)
	if err != nil {
		return &types.LoginResponse{
			Response: types.Response{
				Code:    response.ParameterErrorCode,
				Message: "Failed to generate token",
				Info:    err.Error(),
			},
		}, nil

	}

	// Fetch user's organizations
	modelOrgs, err := l.svcCtx.OrgsUsersModel.FindAllByUserId(l.ctx, user.UserId)
	if err != nil {
		return &types.LoginResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to fetch user organizations",
				Info:    err.Error(),
			},
		}, nil
	}

	// Convert model.Organizations to types.Org
	orgs := make([]types.Org, len(*modelOrgs))
	for i, org := range *modelOrgs {
		orgs[i] = types.Org{
			OrgId:   org.OrgId,
			OrgName: org.OrgName,
			Role:    org.Role,
		}
	}

	// Successful response
	return &types.LoginResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Login successful",
			Info:    "User authenticated successfully",
		},
		Data: &types.LoginResponseData{
			AccessToken:  token,
			AccessExpire: now + accessExpire,
			Orgs:         orgs,
		},
	}, nil
}
