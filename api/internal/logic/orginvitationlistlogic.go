package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgInvitationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgInvitationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgInvitationListLogic {
	return &OrgInvitationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrgInvitationListLogic) OrgInvitationList(req *types.OrgInvitationListReq) (resp *types.OrgInvitationListRsp, err error) {
	// Get user id from context
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	resp = &types.OrgInvitationListRsp{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.OrgInvitationListRspData{
			Invitations: make([]types.OrgInvitation, 0),
		},
	}

	// Get invitations with pagination
	invitations, err := l.svcCtx.InvitationModel.FindAll(l.ctx, uid, int64(req.Page), int64(req.PageSize))
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "Failed to fetch invitations"
		resp.Info = err.Error()
		return resp, nil
	}

	// Convert model to response type
	for _, inv := range *invitations {
		orgName := ""
		org, err := l.svcCtx.OrganizationModel.FindOne(l.ctx, inv.OrgId)
		if err == nil {
			orgName = org.OrgName
		}
		//TODO: TOO SLOW

		resp.Data.Invitations = append(resp.Data.Invitations, types.OrgInvitation{
			OrgId:     inv.OrgId,
			OrgName:   orgName,
			InviterId: inv.InviterId,
			InviteeId: inv.InviteeId,
			Email:     inv.InviteeEmail,
			Role:      inv.Role,
			Status:    inv.Status,
		})
	}

	// Get total count
	total, err := l.svcCtx.InvitationModel.Count(l.ctx, uid)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "Failed to get total count"
		resp.Info = err.Error()
		return resp, nil
	}

	resp.Data.Total = total

	return resp, nil
}
