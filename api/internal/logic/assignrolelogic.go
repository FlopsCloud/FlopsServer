package logic

import (
	"context"
	"database/sql"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRoleLogic {
	return &AssignRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRoleLogic) AssignRole(req *types.AssignRoleRequest) (resp response.Response) {
	// Check if user exists
	_, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.InvalidRequestParamCode, "User not found")
		}
		return response.Error(err.Error())
	}

	// Check if role exists
	_, err = l.svcCtx.RolesModel.FindOne(l.ctx, req.RoleId)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.InvalidRequestParamCode, "Role not found")
		}
		return response.Error(err.Error())
	}

	// Check if user already has this role
	existingUserRole, err := l.svcCtx.UserRolesModel.FindOneByUserIdRoleId(l.ctx, req.UserId, req.RoleId)
	if err != nil && err != model.ErrNotFound {
		return response.Error(err.Error())
	}
	if existingUserRole != nil {
		return response.Fail(response.InvalidRequestParamCode, "User already has this role")
	}

	userRole := &model.UserRoles{
		UserId: sql.NullInt64{
			Int64: int64(req.UserId),
			Valid: true,
		},
		RoleId: sql.NullInt64{
			Int64: int64(req.RoleId),
			Valid: true,
		},
		AssignedAt: time.Now(),
	}

	_, err = l.svcCtx.UserRolesModel.Insert(l.ctx, userRole)
	if err != nil {
		return response.Error(err.Error())
	}

	return response.OK(nil)
}
