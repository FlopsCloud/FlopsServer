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

type BindPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPhoneLogic {
	return &BindPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindPhoneLogic) BindPhone(req *types.BindPhoneRequest) (resp *types.Response, err error) {
	newResponse := func(code int64, message string, info string) *types.Response {
		return &types.Response{
			Code:    code,
			Message: message,
			Info:    info,
		}
	}

	// Verify the provided code
	verificationCode, err := l.svcCtx.VerificationCodesModel.FindOneByPhone(l.ctx, req.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			logx.Error("电话验证码不存在或已过期")
			return newResponse(response.ParameterErrorCode, "电话验证码不存在或已过期", err.Error()), nil
		}
		return newResponse(response.ServerErrorCode, "服务器错误，请联系客服", "(FindOneByPhone)"+err.Error()), nil
	}

	logx.Info("verificationCode", verificationCode)

	if verificationCode.Code != req.Code {
		logx.Error("电话验证码错误", verificationCode.Code, req.Code)
		return newResponse(response.ParameterErrorCode, "电话验证码错误", "The provided code does not match"), nil
	}

	// Check if the code has expired (assuming 5 minutes expiration)
	if time.Now().After(verificationCode.ExpirationTime) {
		logx.Error("电话验证码已过期", verificationCode.ExpirationTime)
		return newResponse(response.ParameterErrorCode, "Verification code expired", "Please request a new verification code"), nil
	}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(uid))

	if err != nil {
		if err == model.ErrNotFound {
			return newResponse(response.ParameterErrorCode, "用户不存在", "User not found"), nil
		} else {
			return newResponse(response.ServerErrorCode, "服务器错误，请联系客服", "(FindOneByPhone)"+err.Error()), nil
		}

	}

	user.Phone = req.Phone
	user.IsPhoneVerified = 1
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		return newResponse(response.ServerErrorCode, "服务器错误，请联系客服", "(Update)"+err.Error()), nil
	}

	// Delete the used verification code
	err = l.svcCtx.VerificationCodesModel.Delete(l.ctx, verificationCode.Id)
	if err != nil {
		logx.Error("Failed to delete verification code:", err)
	}

	return newResponse(response.SuccessCode, "电话号码绑定成功", ""), nil

}
