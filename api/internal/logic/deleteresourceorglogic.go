package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteResourceOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteResourceOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteResourceOrgLogic {
	return &DeleteResourceOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteResourceOrgLogic) DeleteResourceOrg(req *types.DeleteResourceOrgRequest) (resp *types.Response, err error) {
	// Delete the resource-org relationship
	err = l.svcCtx.ResourceOrgsModel.DeleteByUserIDOrgID(l.ctx, req.ResourceId, req.OrgId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Delete Fail",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Delete Success",
	}, nil
}
