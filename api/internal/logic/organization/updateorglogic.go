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

type UpdateOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrgLogic {
	return &UpdateOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrgLogic) UpdateOrg(req *types.UpdateOrgRequest) (resp *types.Response, err error) {

	resp = &types.Response{}

	// 提取用户ID

	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}
	userId := uint64(uid)
	logx.Infof("UpdateOrgLogic UpdateOrg userId %d Orgid %d", userId, req.OrgId)
	//	检查组织是否存在
	existing, err := l.svcCtx.OrganizationModel.FindByUserIdOrgId(l.ctx, userId, req.OrgId)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 400
			resp.Message = "Cannot find your organization"
			resp.Info = err.Error()
			return resp, nil
		} else {
			resp.Code = response.ServerErrorCode
			resp.Message = err.Error()
			return resp, nil
		}

	}

	// 更新组织信息
	existing.OrgName = req.OrgName
	existing.IsPrivate = req.IsPrivate
	existing.IsDefault = req.IsDefault

	err = l.svcCtx.OrganizationModel.Update(l.ctx, existing)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}
	resp.Code = response.SuccessCode
	resp.Message = "Update success"

	return resp, nil
}
