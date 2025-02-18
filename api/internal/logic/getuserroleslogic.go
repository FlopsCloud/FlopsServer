package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRolesLogic {
	return &GetUserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserRolesLogic) GetUserRoles(userId uint64) (resp *types.UserRolesResponse) {
	resp = &types.UserRolesResponse{
		Response: types.Response{},
	}

	// Check if user exists
	_, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Response = types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "User not found",
			}
			return resp
		}
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}

	// Get user roles
	userRoles, err := l.svcCtx.UserRolesModel.FindByUserId(l.ctx, userId)
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}

	var roles []types.Role
	for _, userRole := range userRoles {
		role, err := l.svcCtx.RolesModel.FindOne(l.ctx, uint64(userRole.RoleId.Int64))
		if err != nil {
			continue // Skip if role not found
		}
		roles = append(roles, types.Role{
			Id:          role.RoleId,
			Name:        role.RoleName,
			Description: role.Description.String,
			CreatedAt:   role.CreatedAt.Unix(),
			UpdatedAt:   role.UpdatedAt.Unix(),
		})
	}

	resp.Response = types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}
	resp.Data = types.UserRolesResponseData{
		Roles: roles,
	}

	return resp
}
