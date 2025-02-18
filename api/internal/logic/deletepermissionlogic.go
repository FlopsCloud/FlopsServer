package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeletePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePermissionLogic) DeletePermission(id uint64) (resp response.Response) {
	// Check if permission exists
	permission, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, id)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.InvalidRequestParamCode, "Permission not found")
		}
		return response.Error(err.Error())
	}

	err = sqlx.NewMysql(l.svcCtx.Config.MySQL.DataSource).Transact(func(session sqlx.Session) error {
		// Delete role permissions first
		rolePermissionsModel := l.svcCtx.RolePermissionsModel.WithSession(session)
		err := rolePermissionsModel.DeleteByPermissionId(l.ctx, permission.PermissionId)
		if err != nil {
			return err
		}

		// Delete permission
		permissionsModel := l.svcCtx.PermissionsModel.WithSession(session)
		return permissionsModel.Delete(l.ctx, permission.PermissionId)
	})

	if err != nil {
		return response.Error(err.Error())
	}

	return response.OK(nil)
}
