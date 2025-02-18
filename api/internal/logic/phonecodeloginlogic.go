package logic

import (
	"context"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/jwtx"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type PhonecodeloginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPhonecodeloginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PhonecodeloginLogic {
	return &PhonecodeloginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PhonecodeloginLogic) Phonecodelogin(req *types.PhonecodeloginRequest) (*types.LoginResponse, error) {
	// Helper function to generate error response
	errorResponse := func(code int64, message string, info string) *types.LoginResponse {
		return &types.LoginResponse{
			Response: types.Response{
				Code:    code,
				Message: message,
				Info:    info,
			},
		}
	}

	// Find user by phone number
	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, req.Phonenum)
	if err != nil {
		if err == model.ErrNotFound {
			logx.Error("用户不存在")
			return errorResponse(404, "User not found", "The provided phone number is not associated with any user"), nil
		}
		return errorResponse(500, "Internal server error", "An unexpected error occurred"), nil
	}

	// Verify the provided code
	verificationCode, err := l.svcCtx.VerificationCodesModel.FindOneByPhone(l.ctx, req.Phonenum)
	if err != nil {
		if err == model.ErrNotFound {
			logx.Error("电话验证码不存在或已过期")
			return errorResponse(404, "Verification code not found or expired", "Please request a new verification code"), nil
		}
		return errorResponse(500, "Internal server error", "An unexpected error occurred"), nil
	}

	if verificationCode.Code != req.Code {
		logx.Error("电话验证码错误", verificationCode.Code, req.Code)
		return errorResponse(400, "Invalid verification code", "The provided code does not match"), nil
	}

	// Check if the code has expired (assuming 5 minutes expiration)
	if time.Now().After(verificationCode.ExpirationTime) {
		logx.Error("电话验证码已过期", verificationCode.ExpirationTime)
		return errorResponse(400, "Verification code expired", "Please request a new verification code"), nil
	}

	// Generate JWT token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	token, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, int64(user.UserId), user.Email)
	if err != nil {
		return errorResponse(500, "Failed to generate token", "An unexpected error occurred"), nil
	}

	// Delete the used verification code
	err = l.svcCtx.VerificationCodesModel.Delete(l.ctx, verificationCode.Id)
	if err != nil {
		logx.Error("Failed to delete verification code:", err)
	}

	// Fetch user's organizations
	modelOrgs, err := l.svcCtx.OrgsUsersModel.FindAllByUserId(l.ctx, user.UserId)
	if err != nil {
		return errorResponse(500, "Failed to fetch user organizations", "An unexpected error occurred"), nil
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
			Code:    200,
			Message: "Login successful",
			Info:    "User authenticated successfully",
		},
		Data: types.LoginResponseData{
			AccessToken:  token,
			AccessExpire: now + accessExpire,
			Orgs:         orgs,
		},
	}, nil
}
