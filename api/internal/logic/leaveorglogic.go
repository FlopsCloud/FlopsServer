package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// User leave an organization
func NewLeaveOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveOrgLogic {
	return &LeaveOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LeaveOrgLogic) LeaveOrg(req *types.LeaveOrgReq) (resp *types.Response, err error) {
	if req.OrgId == 0 || req.UserId == 0 {
		return &types.Response{
			Code:    400,
			Message: "Invalid input: OrgId and UserId are required",
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
	l.Logger.Infof("User %s (ID: %d) leave organization %s (ID: %d) ", user.Email, user.UserId, org.OrgName, org.OrgId)

	err = l.svcCtx.OrgsUsersModel.DeleteEx(l.ctx, req.OrgId, req.UserId)

	if err != nil {
		return &types.Response{
			Code:    500,
			Message: "Leav Org fail",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "Leav Org OK",
	}, nil
}
