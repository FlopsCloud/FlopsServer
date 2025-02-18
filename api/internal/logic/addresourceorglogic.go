package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddResourceOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddResourceOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddResourceOrgLogic {
	return &AddResourceOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddResourceOrgLogic) AddResourceOrg(req *types.AddResourceOrgRequest) (resp *types.ResourceOrgResponse, err error) {

	// Create new resource-org relationship
	_, err = l.svcCtx.ResourceOrgsModel.Insert(l.ctx, &model.ResourceOrgs{
		ResourceId: req.ResourceId,
		OrgId:      req.OrgId,
	})
	if err != nil {
		return &types.ResourceOrgResponse{
			Response: types.Response{
				Code:    response.SuccessCode,
				Message: l.svcCtx.W.T("Insert Fail"),
				Info:    err.Error(),
			},
		}, nil
	}

	return &types.ResourceOrgResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ResourceOrg{
			Id:         0,
			ResourceId: req.ResourceId,
			OrgId:      req.OrgId,
		},
	}, nil
}
