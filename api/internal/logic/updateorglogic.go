package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrgLogic {
	return &UpdateOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrgLogic) UpdateOrg(req *types.UpdateOrgRequest) (resp *types.Response, err error) {

	// Call a method on the service context to update the organization

	// Check if user is the org owner
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return &types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "Failed to get user ID",
			Info:    err.Error(),
		}, nil
	}
	userID := uint64(uid)

	org, err := l.svcCtx.OrganizationModel.FindOne(l.ctx, req.OrgId)

	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    response.InvalidRequestParamCode,
				Message: "Organization not found",
			}, nil
		}
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Database error",
			Info:    err.Error(),
		}, nil
	}

	if org.CreatedBy != userID {
		return &types.Response{
			Code:    response.InvalidRequestParamCode,
			Message: "Not authorized - only org owner can update organization",
		}, nil
	}

	org.OrgId = req.OrgId
	org.OrgName = req.OrgName
	org.IsPrivate = req.IsPrivate

	err = l.svcCtx.OrganizationModel.Update(l.ctx, org)
	if err != nil {
		l.Logger.Errorf("Failed to update organization: %v", err)
		return nil, err
	}

	// Create a successful response
	resp = &types.Response{
		Code:    response.SuccessCode,
		Message: "Organization updated successfully",
	}

	// Log the successful update
	l.Logger.Infof("Successfully updated organization ID: %s", org.OrgId)

	return resp, nil
}
