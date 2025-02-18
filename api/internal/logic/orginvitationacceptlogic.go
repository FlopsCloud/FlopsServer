package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgInvitationAcceptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgInvitationAcceptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgInvitationAcceptLogic {
	return &OrgInvitationAcceptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrgInvitationAcceptLogic) OrgInvitationAccept(req *types.OrgInvitationAcceptReq) (resp *types.Response, err error) {
	// Get user id from context
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	userId := uint64(uid)

	// Get user email
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to get user information",
			Info:    err.Error(),
		}, nil
	}

	// Find invitation by token
	invitation, err := l.svcCtx.InvitationModel.FindOne(l.ctx, req.InvitationId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "Invalid or expired invitation token",
			}, nil
		}
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to verify invitation",
			Info:    err.Error(),
		}, nil
	}

	// Verify invitation is for this user
	if invitation.InviteeEmail != user.Email {
		return &types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "This invitation is not for your email address",
			Info:    invitation.InviteeEmail + " " + user.Email,
		}, nil
	}

	if invitation.InviteeId != 0 && invitation.InviteeId != userId {
		return &types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "This invitation if not for you",
			Info:    strconv.FormatUint(invitation.InviteeId, 10) + " " + strconv.FormatUint(userId, 10),
		}, nil
	}
	// Check if invitation is still pending
	if invitation.Status != "pending" {
		return &types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "This invitation has already been used or cancelled",
		}, nil
	}

	// Check if user is already a member of the organization
	_, err = l.svcCtx.OrgsUsersModel.FindOneByOrgIdUserId(l.ctx, invitation.OrgId, userId)
	if err == nil {
		return &types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "You are already a member of this organization",
		}, nil
	}

	// Add user to organization
	_, err = l.svcCtx.OrgsUsersModel.Insert(l.ctx, &model.OrgsUsers{
		OrgId:     invitation.OrgId,
		UserId:    userId,
		Role:      invitation.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to add user to organization",
			Info:    err.Error(),
		}, nil
	}

	// Update invitation status to accepted
	invitation.Status = "accepted"
	err = l.svcCtx.InvitationModel.Update(l.ctx, invitation)
	if err != nil {
		logx.Error("Failed to update invitation status:", err)
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Successfully joined the organization",
	}, nil
}
