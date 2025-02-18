package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListServerOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListServerOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListServerOrgLogic {
	return &ListServerOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListServerOrgLogic) ListServerOrg(req *types.ListServerOrgRequest) (resp *types.ListServerOrgRsp, err error) {
	// todo: add your logic here and delete this line

	serverorgs, err := l.svcCtx.ServerOrgsModel.FindAllByServerOrg(l.ctx, req.OrgId, req.ServerId)
	if err != nil {
		return &types.ListServerOrgRsp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "服务器错误",
				Info:    err.Error(),
			},
		}, nil

	}
	resp = &types.ListServerOrgRsp{}

	for _, serverorg := range serverorgs {
		resp.Data.ServerOrg = append(resp.Data.ServerOrg, types.ServerOrg{
			Id:            uint64(serverorg.Id),
			ServerId:      serverorg.ServerId,
			OrgId:         serverorg.OrgId,
			CreatedAt:     serverorg.CreatedAt.Unix(),
			UpdatedAt:     serverorg.UpdatedAt.Unix(),
			SrvDiscountId: serverorg.SrvDiscountId,
		})
	}
	resp.Data.Total = int64(len(serverorgs))

	return resp, nil
}
