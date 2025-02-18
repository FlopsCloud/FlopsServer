package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(id uint64) (resp response.Response) {
	// Check if role exists
	role, err := l.svcCtx.RolesModel.FindOne(l.ctx, id)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.InvalidRequestParamCode, "Role not found")
		}
		return response.Error(err.Error())
	}

	err = sqlx.NewMysql(l.svcCtx.Config.MySQL.DataSource).Transact(func(session sqlx.Session) error {
		// Delete role permissions first
		rolePermissionsModel := l.svcCtx.RolePermissionsModel.WithSession(session)
		err := rolePermissionsModel.DeleteByRoleId(l.ctx, role.RoleId)
		if err != nil {
			return err
		}

		// Delete user roles
		userRolesModel := l.svcCtx.UserRolesModel.WithSession(session)
		err = userRolesModel.DeleteByRoleId(l.ctx, role.RoleId)
		if err != nil {
			return err
		}

		// Delete role
		rolesModel := l.svcCtx.RolesModel.WithSession(session)
		return rolesModel.Delete(l.ctx, role.RoleId)
	})

	if err != nil {
		return response.Error(err.Error())
	}

	return response.OK(nil)
}
