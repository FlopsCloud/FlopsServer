package logic

import (
	"context"
	"database/sql"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyToJoinOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyToJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyToJoinOrgLogic {
	return &ApplyToJoinOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyToJoinOrgLogic) ApplyToJoinOrg(req *types.ApplyToJoinOrgReq) (resp *types.Response, err error) {
	// Validate input
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

	// Since we don't have an OrgsUsersModel or ApplyJoinModel, we can't directly check if the user is already a member
	// or create an application. Instead, we'll log this information and return a success message.

	// Log the application
	l.Logger.Infof("User %s (ID: %d) applied to join organization %s (ID: %d). Message: %s", user.Email, user.UserId, org.OrgName, org.OrgId, req.Message)

	// In a real implementation, you would typically:
	// 1. Check if the user is already a member of the organization
	// 2. Check if there's an existing pending application
	// 3. Create a new application record in the database

	var nullStr sql.NullString
	if req.Message != "" {
		nullStr = sql.NullString{String: req.Message, Valid: true}
	} else {
		nullStr = sql.NullString{Valid: false} // 表示数据库中的NULL
	}

	_, err = l.svcCtx.ApplyJoinModel.Insert(l.ctx, &model.ApplyJoin{
		UserId:  req.UserId,
		OrgId:   req.OrgId,
		Message: nullStr,
		Status:  "pending",
	})
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: "Invitation fail",
			Info:    err.Error(),
		}, nil
	}

	// if err != nil {
	// 	return &types.Response{
	// 		Code:    500,
	// 		Message: "Invitation sendInvitationEmail fail",
	// 		Info:    err.Error(),
	// 	}, nil
	// }

	return &types.Response{
		Code:    200,
		Message: "Application to join organization logged successfully. Implementation for creating application is pending.",
	}, nil
}
