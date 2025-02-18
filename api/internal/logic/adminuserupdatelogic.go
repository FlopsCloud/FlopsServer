package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserUpdateLogic {
	return &AdminUserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserUpdateLogic) AdminUserUpdate(req *types.AdminUserUpdateRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	// Initialize response
	resp = &types.Response{
		Code:    response.SuccessCode,
		Message: "success",
	}

	role, _ := l.ctx.Value("role").(string)

	if req.SysRole == "superadmin" && role != "superadmin" {
		resp.Code = response.UnauthorizedCode
		resp.Message = "permission denied"
		resp.Info = role
		return resp, nil
	}

	if role != "superadmin" && role != "admin" {
		resp.Code = response.UnauthorizedCode
		resp.Message = "permission denied"
		resp.Info = role
		return resp, nil
	}

	if req.SysRole != "" && (req.SysRole != "superadmin" && req.SysRole != "admin" && req.SysRole != "user" && req.SysRole != "member") {
		resp.Code = response.ParameterErrorCode
		resp.Message = "invalid sysRole, only superadmin, admin, user, member are allowed"
		resp.Info = req.SysRole
		return resp, nil
	}

	// Validate user ID exists
	if req.UserId == 0 {
		resp.Code = response.ParameterErrorCode
		resp.Message = "User ID is required"
		return resp, nil
	}

	// Get existing user
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		resp.Code = response.UserNotExistCode
		resp.Message = "Failed to find user"
		return resp, err
	}

	if req.SysRole != "" {
		user.SysRole = req.SysRole
	}

	// Update fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if req.IsDeleted != 0 {
		user.IsDeleted = req.IsDeleted
	}

	// Save updates
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		resp.Code = 500
		resp.Message = "Failed to update user"
		return resp, err
	}

	return
}
