package admin

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteOrgLogic {
	return &AdminDeleteOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteOrgLogic) AdminDeleteOrg(req *types.DeleteOrgRequest) (resp *types.Response, err error) {
	role, _ := l.ctx.Value("role").(string)
	if role != "superadmin" {
		resp = &types.Response{
			Code:    response.UnauthorizedCode,
			Message: "Permission denied, Super Admin only",
		}
		return resp, nil
	}

	// Check if org exists

	_, err = l.svcCtx.OrganizationModel.FindOne(l.ctx, req.OrgId)
	if err == model.ErrNotFound {
		return &types.Response{
			Code:    response.NotFoundCode,
			Message: "Organization not found",
		}, nil
	}
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}, nil
	}

	// Delete organization
	err = l.svcCtx.OrganizationModel.Delete(l.ctx, req.OrgId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Organization deleted successfully",
	}, nil
}
