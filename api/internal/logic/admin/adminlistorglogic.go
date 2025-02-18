package admin

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListOrgLogic {
	return &AdminListOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListOrgLogic) AdminListOrg(req *types.ListOrgRequest) (resp *types.ListOrgResp, err error) {

	role, _ := l.ctx.Value("role").(string)
	if role != "superadmin" {
		resp = &types.ListOrgResp{
			Response: types.Response{
				Code:    response.UnauthorizedCode,
				Message: "Permission denied, Super Admin only",
				Info:    role,
			},
		}
		return resp, nil
	}

	// Get all organizations
	orgs, err := l.svcCtx.OrganizationModel.FindAllEx(l.ctx)
	if err != nil {
		return &types.ListOrgResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: err.Error(),
			},
		}, nil
	}

	// Convert to response format
	var orgList []types.Organization
	for _, org := range *orgs {
		orgList = append(orgList, types.Organization{
			OrgId:     org.OrgId,
			OrgName:   org.OrgName,
			IsPrivate: org.IsPrivate,
			CreatedBy: org.CreatedBy,
			UpdatedAt: uint64(org.UpdatedAt.Unix()),
			CreatedAt: uint64(org.CreatedAt.Unix()),
			Username:  org.Username,
			IsDefault: uint64(org.IsDefault),
		})
	}

	return &types.ListOrgResp{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Organizations retrieved successfully",
		},
		Data: types.ListOrgRespData{
			Orgs: orgList,
		},
	}, nil

}
