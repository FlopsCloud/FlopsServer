package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserToOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUserToOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserToOrgLogic {
	return &AddUserToOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUserToOrgLogic) AddUserToOrg(req *types.AddUserToOrgReq) (resp *types.Response, err error) {
	// Validate input
	if req.OrgId == 0 || req.UserId == 0 || req.Role == "" {
		return &types.Response{
			Code:    400,
			Message: "Invalid input: OrgId, UserId, and Role are required",
		}, nil
	}

	// Validate role
	if req.Role != "member" && req.Role != "guest" && req.Role != "manager" && req.Role != "admin" {
		return &types.Response{
			Code:    400,
			Message: "Invalid role: must be either 'member', 'guest', 'manager', or 'admin'",
		}, nil
	}

	// Check if the organization exists
	org, err := l.svcCtx.OrganizationModel.FindOne(l.ctx, req.OrgId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    404,
				Message: "Organization not found",
			}, nil
		}
		return nil, err
	}

	// Check if the user exists
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    404,
				Message: "User not found",
			}, nil
		}
		return nil, err
	}

	// Since we don't have an OrgsUsersModel, we can't directly check if the user is already a member
	// or add them to the organization. Instead, we'll log this information and return a success message.

	// Log the action
	l.Logger.Infof("Adding user %s (ID: %d) to organization %s (ID: %d) with role %s", user.Email, user.UserId, org.OrgName, org.OrgId, req.Role)

	// In a real implementation, you would typically:
	// 1. Check if the user is already a member of the organization
	// 2. Add the user to the organization with the specified role
	// 3. Update any necessary records or permissions

	_, err = l.svcCtx.OrgsUsersModel.Insert(l.ctx, &model.OrgsUsers{
		UserId: req.UserId,
		OrgId:  req.OrgId,
		Role:   req.Role,
	})

	if err != nil {
		return &types.Response{
			Code:    500,
			Message: "add user to org fail",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "User addition to organization successfully ",
	}, nil
}
