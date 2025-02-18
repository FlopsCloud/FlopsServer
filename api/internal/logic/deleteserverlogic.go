package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteServerLogic {
	return &DeleteServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteServerLogic) DeleteServer(req *types.DeleteServerReq) (resp *types.Response, err error) {
	sysrole, _ := l.ctx.Value("role").(string)
	if sysrole != "admin" && sysrole != "superadmin" {
		var res types.Response
		res.Code = response.UnauthorizedCode
		res.Message = "only admin can access"
		return &res, nil

	}

	// servers, err := l.svcCtx.ServerOrgsModel.FindAllByServerOrg(l.ctx, 0, req.ServerId)
	// if len(servers) == 1 && err == nil {
	// 	// Delete server-org relationship
	// 	err = l.svcCtx.ServerOrgsModel.DeleteByServerOrg(l.ctx, &model.ServerOrgs{
	// 		OrgId:    servers[0].OrgId,
	// 		ServerId: req.ServerId,
	// 	})
	// 	if err != nil {
	// 		return &types.Response{
	// 			Code:    response.ServerErrorCode,
	// 			Message: "ServerOrg Delete Failed",
	// 			Info:    err.Error(),
	// 		}, nil
	// 	}
	// } else if len(servers) > 1 {
	// 	return &types.Response{
	// 		Code:    response.ServerErrorCode,
	// 		Message: "more than one server-org relationship",
	// 		Info:    "more than one server-org relationship",
	// 	}, nil
	// } else if err != nil {
	// 	return &types.Response{
	// 		Code:    response.ServerErrorCode,
	// 		Message: "ServerOrg Delete Failed",
	// 		Info:    err.Error(),
	// 	}, nil
	// }

	// Delete server
	err = l.svcCtx.ServerModel.Delete(l.ctx, req.ServerId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Server Delete Failed",
			Info:    err.Error(),
		}, nil
	}
	return &types.Response{
		Code:    response.SuccessCode,
		Message: "ok",
	}, nil
}
