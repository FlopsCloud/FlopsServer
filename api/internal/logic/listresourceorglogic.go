package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListResourceOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListResourceOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListResourceOrgLogic {
	return &ListResourceOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListResourceOrgLogic) ListResourceOrg(req *types.ListResourceOrgRequest) (resp *types.ListResourceOrgResp, err error) {
	// Check admin permissions
	role, _ := l.ctx.Value("role").(string)
	if role != "admin" && role != "superadmin" {
		return &types.ListResourceOrgResp{
			Response: types.Response{
				Code:    response.UnauthorizedCode,
				Message: "Permission denied, Admin only",
			},
		}, nil
	}

	var resourceOrgs []*model.ResourceOrgs
	var total int64
	var err1 error

	// Get resource-org mappings based on filters
	if req.ResourceId > 0 && req.OrgId > 0 {
		// Filter by both ResourceId and OrgId
		result, err := l.svcCtx.ResourceOrgsModel.FindByResourceIdOrgId(l.ctx, req.ResourceId, req.OrgId)
		if err != nil {
			total = 0
			resourceOrgs = []*model.ResourceOrgs{}
		} else {
			total = 1
			resourceOrgs = []*model.ResourceOrgs{result}
		}
	} else if req.ResourceId > 0 {
		// Filter by ResourceId
		resourceOrgs, err1 = l.svcCtx.ResourceOrgsModel.FindByResourceId(l.ctx, req.ResourceId)
		if err1 != nil {
			return &types.ListResourceOrgResp{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: err1.Error(),
				},
			}, nil
		}
		total = int64(len(resourceOrgs))
	} else if req.OrgId > 0 {
		// Filter by OrgId
		resourceOrgs, err1 = l.svcCtx.ResourceOrgsModel.FindByOrgId(l.ctx, req.OrgId)
		if err1 != nil {
			return &types.ListResourceOrgResp{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: err1.Error(),
				},
			}, nil
		}
		total = int64(len(resourceOrgs))
	} else {
		// Get all resource-org mappings
		resourceOrgs, err1 = l.svcCtx.ResourceOrgsModel.FindAll(l.ctx)
		if err1 != nil {
			return &types.ListResourceOrgResp{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: err1.Error(),
				},
			}, nil
		}
		total = int64(len(resourceOrgs))
	}

	// Convert model to response type
	var result []types.ResourceOrg
	for _, ro := range resourceOrgs {
		result = append(result, types.ResourceOrg{
			Id:         uint64(ro.Id),
			ResourceId: ro.ResourceId,
			OrgId:      ro.OrgId,
			DiscountId: ro.DiscountId,

			CreatedAt: int64(ro.CreatedAt.Unix()),
			UpdatedAt: int64(ro.UpdatedAt.Unix()),
		})
	}

	return &types.ListResourceOrgResp{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ListResourceOrgRespData{
			Resources: result,
			Total:     uint64(total),
		},
	}, nil
}
