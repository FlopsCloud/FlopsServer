package logic

import (
	"context"
	"encoding/json"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindEmailLogic {
	return &BindEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindEmailLogic) BindEmail(req *types.BindEmailRequest) (resp *types.Response, err error) {
	newResponse := func(code int64, message string, info string) *types.Response {
		return &types.Response{
			Code:    code,
			Message: message,
			Info:    info,
		}
	}

	// Verify the provided code
	verificationCode, err := l.svcCtx.VerificationCodesModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil {
		if err == model.ErrNotFound {
			logx.Error("邮箱验证码不存在或已过期")
			return newResponse(response.ParameterErrorCode, "邮箱验证码不存在或已过期", err.Error()), nil
		}
		return newResponse(response.ServerErrorCode, "服务器错误，请联系客服", "(FindOneByEmail)"+err.Error()), nil
	}

	logx.Info("verificationCode", verificationCode)

	if verificationCode.Code != req.Code {
		logx.Error("邮箱验证码错误", verificationCode.Code, req.Code)
		return newResponse(response.ParameterErrorCode, "邮箱验证码错误", "The provided code does not match"), nil
	}

	// Check if the code has expired (assuming 5 minutes expiration)
	if time.Now().After(verificationCode.ExpirationTime) {
		logx.Error("邮箱验证码已过期", verificationCode.ExpirationTime)
		return newResponse(response.ParameterErrorCode, "邮箱验证码已过期", "Please request a new verification code"), nil
	}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(uid))
	if err != nil {
		if err == model.ErrNotFound {
			return newResponse(response.ParameterErrorCode, "用户不存在", "User not found"), nil
		} else {
			return newResponse(response.ServerErrorCode, "服务器错误，请联系客服", "(FindOneByEmail)"+err.Error()), nil
		}

	}

	user.Email = req.Email
	user.IsEmailVerified = 1
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		return newResponse(response.ServerErrorCode, "服务器错误，请联系客服", "(Update)"+err.Error()), nil
	}

	// Delete the used verification code
	err = l.svcCtx.VerificationCodesModel.Delete(l.ctx, verificationCode.Id)
	if err != nil {
		logx.Error("Failed to delete verification code:", err)
	}

	return newResponse(response.SuccessCode, "绑定成功", ""), nil
}
