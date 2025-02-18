package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type KeepAliveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KeepAliveLogic {
	return &KeepAliveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KeepAliveLogic) KeepAlive() (resp response.Response, err error) {

	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return response.Fail(response.UnauthorizedCode, "Invalid user ID"), nil
	}
	userId := uint64(uid)

	// Create a response with the userId
	return response.OK(map[string]interface{}{
		"message": "Keep-alive successful",
		"userId":  userId,
	}), nil
}
