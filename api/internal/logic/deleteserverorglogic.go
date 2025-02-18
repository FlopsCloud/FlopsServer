package logic

import (
	"context"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteServerOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteServerOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteServerOrgLogic {
	return &DeleteServerOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteServerOrgLogic) DeleteServerOrg(req *types.DeleteServerOrgRequest) (resp *types.Response, err error) {
	// Create a ServerOrgs object with the request data
	serverOrg := &model.ServerOrgs{
		ServerId: req.ServerId,
		OrgId:    req.OrgId,
	}

	// Try to find and delete the record
	err = l.svcCtx.ServerOrgsModel.DeleteByServerOrg(l.ctx, serverOrg)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    404,
				Message: fmt.Sprintf("Server-org association not found for ServerId: %d and OrgId: %d", req.ServerId, req.OrgId),
			}, nil
		}
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "Success",
	}, nil
}
