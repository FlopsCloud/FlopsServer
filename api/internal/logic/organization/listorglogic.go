package organization

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrgLogic {
	return &ListOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOrgLogic) ListOrg(req *types.ListOrgRequest) (resp *types.ListOrgResp, err error) {

	resp = &types.ListOrgResp{}

	// 提取用户ID
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}
	userId := uint64(uid)

	all, err := l.svcCtx.OrganizationModel.FindAllPublicEx(l.ctx, userId)
	if err != nil && err != model.ErrNotFound {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}

	var List []types.Organization
	for _, item := range *all {
		List = append(List, types.Organization{
			OrgId:     item.OrgId,
			OrgName:   item.OrgName,
			CreatedBy: item.CreatedBy,
			UpdatedAt: uint64(item.UpdatedAt.Unix()),
			CreatedAt: uint64(item.CreatedAt.Unix()),
			IsPrivate: item.IsPrivate,
			Username:  item.Username, // 添加用户名字段
			IsDefault: uint64(item.IsDefault),
		})
	}

	resp.Code = response.SuccessCode
	resp.Message = "ok"
	resp.Data.Orgs = List
	resp.Data.Total = int(len(List))

	return resp, nil

}
