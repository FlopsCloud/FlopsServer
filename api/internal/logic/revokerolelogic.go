package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeRoleLogic {
	return &RevokeRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeRoleLogic) RevokeRole(userId, roleId uint64) (resp response.Response) {
	// Check if user exists
	_, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.InvalidRequestParamCode, "User not found")
		}
		return response.Error(err.Error())
	}

	// Check if role exists
	_, err = l.svcCtx.RolesModel.FindOne(l.ctx, roleId)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.InvalidRequestParamCode, "Role not found")
		}
		return response.Error(err.Error())
	}

	// Check if user has this role
	userRole, err := l.svcCtx.UserRolesModel.FindOneByUserIdRoleId(l.ctx, userId, roleId)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.InvalidRequestParamCode, "User does not have this role")
		}
		return response.Error(err.Error())
	}

	err = l.svcCtx.UserRolesModel.Delete(l.ctx, userRole.Id)
	if err != nil {
		return response.Error(err.Error())
	}

	return response.OK(nil)
}
