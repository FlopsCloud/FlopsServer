package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminServerPriceUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminServerPriceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminServerPriceUpdateLogic {
	return &AdminServerPriceUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminServerPriceUpdateLogic) AdminServerPriceUpdate(req *types.AdminServerPriceUpdateRequest) (resp *types.Response, err error) {
	logx.Infof("AdminServerPriceUpdateLogic AdminServerPriceUpdate req: %+v", req)

	// Check admin permissions
	sysrole, _ := l.ctx.Value("role").(string)
	if sysrole != "admin" && sysrole != "superadmin" {
		return &types.Response{
			Code:    response.UnauthorizedCode,
			Message: "only admin can access",
		}, nil
	}

	// Validate input
	// if req.Id <= 0 {
	// 	return &types.Response{
	// 		Code:    response.InvalidRequestParamCode,
	// 		Message: "invalid server id",
	// 	}, nil
	// }

	for _, serverId := range req.ServerIds {
		// Check if server exists
		server, err := l.svcCtx.ServerModel.FindOne(l.ctx, serverId)
		if err != nil {
			if err == model.ErrNotFound {
				return &types.Response{
					Code:    response.NotFoundCode,
					Message: "server not found",
				}, nil
			}
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to get server",
				Info:    err.Error(),
			}, nil
		}

		// Update server pay prices
		server.PayPrices = req.PayPrices
		err = l.svcCtx.ServerModel.Update(l.ctx, server)
		if err != nil {
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: "failed to update server pay prices",
				Info:    err.Error(),
			}, nil
		}
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "success",
	}, nil
}
