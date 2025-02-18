package logic

import (
	"context"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateServerOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateServerOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateServerOrgLogic {
	return &UpdateServerOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateServerOrgLogic) UpdateServerOrg(req *types.UpdateServerOrgRequest) (resp *types.Response, err error) {

	//TODO  : if no org , create a new one

	serverOrg, err := l.svcCtx.ServerOrgsModel.FindOneByServerOrg(l.ctx, req.OrgId, req.ServerId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "UpdateServergo Org Fail",
			Info:    err.Error(),
		}, nil
	}

	serverOrg.SrvDiscountId = req.SrvDiscountId

	serverOrg.UpdatedAt = time.Now()

	err = l.svcCtx.ServerOrgsModel.Update(l.ctx, serverOrg)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "UpdateServergo Org Fail",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}, nil

}
