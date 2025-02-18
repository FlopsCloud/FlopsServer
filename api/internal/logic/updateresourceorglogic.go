package logic

import (
	"context"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateResourceOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateResourceOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateResourceOrgLogic {
	return &UpdateResourceOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateResourceOrgLogic) UpdateResourceOrg(req *types.UpdateResourceOrgRequest) (resp *types.Response, err error) {
	// Check admin permissions
	sysrole, _ := l.ctx.Value("role").(string)
	if sysrole != "admin" && sysrole != "superadmin" {
		return &types.Response{
			Code:    response.UnauthorizedCode,
			Message: "only admin can access",
		}, nil
	}

	// Check if resource-org mapping exists
	existing, err := l.svcCtx.ResourceOrgsModel.FindByResourceIdOrgId(l.ctx, req.ResourceId, req.OrgId)
	if err == model.ErrNotFound {
		return &types.Response{
			Code:    response.NotFoundCode,
			Message: "Resource-Organization mapping not found",
		}, nil
	}
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}, nil
	}

	// Update resource-org mapping
	err = l.svcCtx.ResourceOrgsModel.Update(l.ctx, &model.ResourceOrgs{
		Id:         existing.Id,
		ResourceId: req.ResourceId,
		OrgId:      req.OrgId,
		DiscountId: req.DiscountId,
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Resource-Organization mapping updated successfully",
	}, nil
}
