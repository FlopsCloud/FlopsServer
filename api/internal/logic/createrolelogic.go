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
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleRequest) (resp *types.RoleResponse) {
	resp = &types.RoleResponse{
		Response: types.Response{},
	}

	// Check if role name already exists
	existingRole, err := l.svcCtx.RolesModel.FindOneByRoleName(l.ctx, req.Name)
	if err != nil && err != model.ErrNotFound {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}
	if existingRole != nil {
		resp.Response = types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "Role name already exists",
		}
		return resp
	}

	role := &model.Roles{
		RoleName: req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = sqlx.NewMysql(l.svcCtx.Config.MySQL.DataSource).Transact(func(session sqlx.Session) error {
		// Create role with transaction
		rolesModel := l.svcCtx.RolesModel.WithSession(session)
		result, err := rolesModel.Insert(l.ctx, role)
		if err != nil {
			return err
		}

		roleId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Insert role permissions
		if len(req.Permissions) > 0 {
			rolePermissionsModel := l.svcCtx.RolePermissionsModel.WithSession(session)
			for _, permissionId := range req.Permissions {
				_, err = rolePermissionsModel.Insert(l.ctx, &model.RolePermissions{
					RoleId:       uint64(roleId),
					PermissionId: uint64(permissionId),
					GrantedAt:    time.Now(),
				})
				if err != nil {
					return err
				}
			}
		}

		role.RoleId = uint64(roleId)
		return nil
	})

	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}

	resp.Response = types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}
	resp.Data = types.Role{
		Id:          role.RoleId,
		Name:        role.RoleName,
		Description: role.Description.String,
		CreatedAt:   role.CreatedAt.Unix(),
		UpdatedAt:   role.UpdatedAt.Unix(),
	}

	return resp
}
