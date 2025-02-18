package admin

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateOrgLogic {
	return &AdminUpdateOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateOrgLogic) AdminUpdateOrg(req *types.UpdateOrgRequest) (resp *types.Response, err error) {
	role, _ := l.ctx.Value("role").(string)
	if role != "superadmin" {
		resp = &types.Response{
			Code:    response.UnauthorizedCode,
			Message: "Permission denied, Super Admin only",
		}
		return resp, nil
	}

	// Check if org exists
	existing, err := l.svcCtx.OrganizationModel.FindOne(l.ctx, req.OrgId)
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

	// Check if new name already exists (if name is being changed)
	if existing.OrgName != req.OrgName {
		exists, err := l.svcCtx.OrganizationModel.FindByOrgName(l.ctx, req.OrgName)
		if err != nil && err != model.ErrNotFound {
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: err.Error(),
			}, nil
		}
		if exists != nil {
			return &types.Response{
				Code:    response.ParameterErrorCode,
				Message: "Organization name already exists",
			}, nil
		}
	}

	// Update organization
	err = l.svcCtx.OrganizationModel.Update(l.ctx, &model.Organizations{
		OrgId:   req.OrgId,
		OrgName: req.OrgName,
		// Description: req.Description,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Organization updated successfully",
	}, nil
}
