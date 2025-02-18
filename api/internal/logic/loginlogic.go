package logic

import (
	"context"
	"fca/common/response"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/cryptx"
	"fca/common/jwtx"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (res response.Response) {

	if !VerifyCaptcha(req.CaptchaId, req.Captcha) {
		logx.Infof("验证码错误: %s, %s", req.CaptchaId, req.Captcha)
		return response.Fail(response.WrongCaptchaCode, "验证码错误")
	}

	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(100, "用户不存在")
		}
		return response.Error(err.Error())
	}

	password := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, req.Password)
	if user.PasswordHash != password {
		return response.Fail(100, "密码错误")
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	token, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, int64(user.UserId), user.Email)
	if err != nil {
		return response.Error(err.Error())
	}

	return response.OK(&types.LoginResponseData{
		AccessToken:  token,
		AccessExpire: now + accessExpire,
	})
}
