package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserDeleteLogic {
	return &AdminUserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserDeleteLogic) AdminUserDelete(req *types.AdminUserDeleteRequest) (resp *types.Response, err error) {
	resp = &types.Response{
		Code:    0,
		Message: "Success",
	}

	// 管理员权限验证 start
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	email, _ := l.ctx.Value("email").(string)

	logx.Info("JWT uid=", uid, " Name=", email)
	if uid != 1 {
		resp.Code = 1
		resp.Message = "Permission denied"
		return resp, nil
	}
	// 管理员权限验证 end

	var user *model.Users

	// Check if both UserId and Email are provided
	if req.UserId > 0 && req.Email != "" {
		user, err = l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
		if err == nil && user.Email != req.Email {
			resp.Code = 1
			resp.Message = "UserId and Email do not match"
			return resp, nil
		}
	} else if req.UserId > 0 {
		user, err = l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	} else if req.Email != "" {
		user, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	} else {
		resp.Code = 1
		resp.Message = "Either UserId or Email must be provided"
		return resp, nil
	}

	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 1
			resp.Message = "User not found"
			return resp, nil
		}
		resp.Code = 1
		resp.Message = "Failed to find user: " + err.Error()
		return resp, err
	}

	err = l.svcCtx.UserModel.Delete(l.ctx, user.UserId)
	if err != nil {
		resp.Code = 1
		resp.Message = "Failed to delete user: " + err.Error()
		return resp, err
	}

	return resp, nil
}
