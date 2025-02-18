package organization

import (
	"context"
	"encoding/json"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrgLogic {
	return &CreateOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrgLogic) CreateOrg(req *types.CreateOrgRequest) (resp *types.CreateOrgResp, err error) {

	resp = &types.CreateOrgResp{}

	// 提取用户ID

	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		resp.Response = types.Response{
			Code:    response.UnauthorizedCode,
			Message: "Can not get user id",
			Info:    err.Error(),
		}
		return resp, nil
	}
	userId := uint64(uid)

	// 检查OrgName是否已存在

	existing, err := l.svcCtx.OrganizationModel.FindByOrgName(l.ctx, req.OrgName)
	if err != nil && err != model.ErrNotFound {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp, nil
	}
	if existing != nil {
		resp.Response = types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "OrgName already exists",
		}
		return resp, nil
	}

	// 插入新组织
	org := &model.Organizations{
		OrgName:   req.OrgName,
		IsPrivate: req.IsPrivate,
		IsDefault: req.IsDefault,
		CreatedBy: userId,
		CreatedAt: time.Now(),
	}

	result, err := l.svcCtx.OrganizationModel.Insert(l.ctx, org)
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp, nil
	}

	// 获取插入的记录ID
	id, err := result.LastInsertId()
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp, nil
	}

	org.OrgId = uint64(id)

	// 插入组织成员
	orgMember := &model.OrgsUsers{
		OrgId:     org.OrgId,
		UserId:    userId,
		Role:      "owner", //owner,guest,member
		CreatedAt: time.Now(),
	}
	_, err = l.svcCtx.OrgsUsersModel.Insert(l.ctx, orgMember)
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp, nil
	}
	// 组织创建成功，返回组织信息

	resp.Response = types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}
	resp.Data = types.Organization{
		OrgId:     org.OrgId,
		OrgName:   org.OrgName,
		CreatedBy: userId,
	}

	return resp, nil
}
