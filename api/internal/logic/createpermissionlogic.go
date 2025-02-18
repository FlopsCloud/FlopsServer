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

type CreatePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePermissionLogic {
	return &CreatePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePermissionLogic) CreatePermission(req *types.CreatePermissionRequest) (resp *types.PermissionResponse) {
	resp = &types.PermissionResponse{
		Response: types.Response{},
	}

	// Check if permission code already exists
	existingPerm, err := l.svcCtx.PermissionsModel.FindOneByCode(l.ctx, req.Code)
	if err != nil && err != model.ErrNotFound {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}
	if existingPerm != nil {
		resp.Response = types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "Permission code already exists",
		}
		return resp
	}

	permission := &model.Permissions{
		PermissionName: req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
		CreatedAt: time.Now(),
	}

	result, err := l.svcCtx.PermissionsModel.Insert(l.ctx, permission)
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}

	id, err := result.LastInsertId()
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}

	permission.PermissionId = uint64(id)

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
