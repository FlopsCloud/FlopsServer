package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyToJoinOrglistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// User applies to join an organization
func NewApplyToJoinOrglistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyToJoinOrglistLogic {
	return &ApplyToJoinOrglistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyToJoinOrglistLogic) ApplyToJoinOrglist(req *types.ApplyToJoinListReq) (resp *types.ApplyToJoinListRsp, err error) {
	// get user id from contex
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	// Initialize response
	resp = &types.ApplyToJoinListRsp{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ApplyToJoinListRspData{
			Applies: make([]types.ApplyRec, 0),
		},
	}

	// Get all applications with pagination
	// Note: You may need to add FindAll method to ApplyJoinModel if not exists
	applications, err := l.svcCtx.ApplyJoinModel.FindAll(l.ctx, uid, int64(req.Page), int64(req.PageSize))
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "Failed to fetch applications"
		resp.Info = err.Error()
		return resp, nil
	}

	// Convert model.ApplyJoin to types.ApplyRec
	for _, app := range *applications {
		message := ""
		if app.Message.Valid {
			message = app.Message.String
		}
		orgName := ""

		org, err := l.svcCtx.OrganizationModel.FindOne(l.ctx, app.OrgId)
		if err != nil {
			orgName = "" // or some default value
		} else {
			orgName = org.OrgName
		}

		resp.Data.Applies = append(resp.Data.Applies, types.ApplyRec{
			OrgId:   app.OrgId,
			UserId:  app.UserId,
			OrgName: orgName,

			Message: message,
		})
	}

	// Get total count
	total, err := l.svcCtx.ApplyJoinModel.Count(l.ctx, uid)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "Failed to get total count"
		resp.Info = err.Error()
		return resp, nil
	}

	resp.Data.Total = total

	return resp, nil
}
