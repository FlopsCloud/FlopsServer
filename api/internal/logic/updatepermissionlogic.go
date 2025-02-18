package logic

import (
	"context"
	"database/sql"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePermissionLogic {
	return &UpdatePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePermissionLogic) UpdatePermission(req *types.UpdatePermissionRequest) (resp *types.PermissionResponse) {
	resp = &types.PermissionResponse{
		Response: types.Response{},
	}

	// Check if permission exists
	permission, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Response = types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "Permission not found",
			}
			return resp
		}
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}

	if req.Name != "" {
		permission.PermissionName = req.Name
	}

	if req.Description != "" {
		permission.Description = sql.NullString{
			String: req.Description,
			Valid:  true,
		}
	}

	err = l.svcCtx.PermissionsModel.Update(l.ctx, permission)
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
	resp.Data = types.Permission{
		Id:          permission.PermissionId,
		Name:        permission.PermissionName,
		Description: permission.Description.String,
		CreatedAt:   permission.CreatedAt.Unix(),
	}

	return resp
}
