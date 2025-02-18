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

type DeleteOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrgLogic {
	return &DeleteOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOrgLogic) DeleteOrg(req *types.DeleteOrgRequest) (resp *types.Response, err error) {
	resp = &types.Response{}

	// 提取用户ID

	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {

		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()

		return resp, nil
	}
	userId := uint64(uid)

	// one, err := l.svcCtx.OrganizationModel.FindByUserIdOrgId(l.ctx, userId, req.OrgId)
	one, err := l.svcCtx.OrgsUsersModel.FindOneByOrgIdUserId(l.ctx, req.OrgId, userId)
	if err != nil && err != model.ErrNotFound {
		resp.Code = response.ServerErrorCode
		resp.Message = err.Error()
		return resp, nil
	}
	if err == model.ErrNotFound {
		resp.Code = response.InvalidRequestParamCode
		resp.Message = "Your are not member of this organization"
		return resp, nil
	}
	if one == nil {
		resp.Code = response.InvalidRequestParamCode
		resp.Message = "Your are not member of this organization"
		return resp, nil
	}
	if one.Role != "owner" {
		resp.Code = response.InvalidRequestParamCode
		resp.Message = "Only owner can delete the organization"
		return resp, nil
	}

	// 检查Org User Name是否已存在

	existing, err := l.svcCtx.OrgsUsersModel.FindAllByOrgId(l.ctx, req.OrgId)
	if err != nil && err != model.ErrNotFound {
		resp.Code = response.ServerErrorCode
		resp.Message = "Cannot find Members of this organization"
		resp.Info = err.Error()
		return resp, nil
	}
	if existing != nil {

		for _, item := range *existing {
			if item.UserId != userId {
				resp.Code = response.InvalidRequestParamCode
				resp.Message = "还有人在组织中，无法删除"
				return resp, nil
			}

		}

	}
	// 删除组织 如果有引用，通常会失败
	err = l.svcCtx.OrganizationModel.Delete(l.ctx, req.OrgId)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "删除组织失败"
		resp.Info = err.Error()

		return resp, nil
	}

	// 删除组织用户关系
	err = l.svcCtx.OrgsUsersModel.DeleteExOwner(l.ctx, userId, req.OrgId)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "删除组织用户关系失败"
		resp.Info = err.Error()
		return resp, nil
	}

	resp.Code = response.SuccessCode
	resp.Message = "success"

	return resp, nil
}
