package logic

import (
	"context"
	"encoding/json"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoRequest) (resp *types.Response, err error) {
	resp = new(types.Response)

	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		resp.Code = response.UnauthorizedCode
		resp.Message = "Invalid user ID"
		return resp, nil
	}
	userId := uint64(uid)

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Message = "User not found"
			return resp, nil
		}
		return resp, err
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	err = l.svcCtx.UserModel.UpdateInfo(l.ctx, &model.Users{
		UserId:   userId,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Nickname: user.Nickname,
	})
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "更新用户信息失败"
		resp.Info = err.Error()
		return resp, nil

	}

	resp.Code = 200
	resp.Message = "success"
	return resp, nil
}
