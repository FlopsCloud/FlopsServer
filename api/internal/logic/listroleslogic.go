package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRolesLogic) ListRoles(req *types.RoleListRequest) (resp *types.RoleListResponse) {
	resp = &types.RoleListResponse{
		Response: types.Response{},
	}

	roles, total, err := l.svcCtx.RolesModel.FindByName(l.ctx, req.Name, req.Page, req.PageSize)
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp
	}

	var roleList []types.Role
	for _, role := range roles {
		roleList = append(roleList, types.Role{
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
	resp.Data = types.RoleListResponseData{
		Roles: roleList,
		Total: total,
	}

	return resp
}
