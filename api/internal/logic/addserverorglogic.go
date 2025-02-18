package logic

import (
	"context"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddServerOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddServerOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddServerOrgLogic {
	return &AddServerOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddServerOrgLogic) AddServerOrg(req *types.AddServerOrgRequest) (resp *types.ServerOrgResponse, err error) {
	serverOrg := &model.ServerOrgs{
		ServerId:  req.ServerId,
		OrgId:     req.OrgId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := l.svcCtx.ServerOrgsModel.Insert(l.ctx, serverOrg)
	if err != nil {
		return &types.ServerOrgResponse{
			Response: types.Response{
				Code:    500,
				Message: "AddServerOrg Fail",
				Info:    err.Error(),
			},
		}, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &types.ServerOrgResponse{
		Response: types.Response{
			Code:    0,
			Message: "Success",
		},
		Data: types.ServerOrg{
			Id:          uint64(id),
			ServerId:    req.ServerId,
			OrgId:       req.OrgId,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
			Description: req.Description,
		},
	}, nil
}
