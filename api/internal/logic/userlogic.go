package logic

import (
	"context"
	"encoding/json"
	"fca/common/response"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) User() (resp response.Response) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(uid))
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.UserNotExistCode, "用户不存在")
		}
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	// Fetch user's organizations
	modelOrgs, err := l.svcCtx.OrgsUsersModel.FindAllByUserId(l.ctx, user.UserId)
	if err != nil {
		resp.Code = 500
		resp.Message = "Failed to fetch user organizations"
		return resp
	}

	// Convert model.Organizations to types.Org
	orgs := make([]types.Org, len(*modelOrgs))
	for i, org := range *modelOrgs {
		orgs[i] = types.Org{
			OrgId:   org.OrgId,
			OrgName: org.OrgName,
			Role:    org.Role,
		}
	}

	return response.OK(&types.UserResponseData{
		UserId:    user.UserId,
		Username:  user.Username,
		Email:     user.Email,
		Balance:   user.Balance,
		Phone:     user.Phone,
		Nickname:  user.Nickname,
		IsDeleted: int64(user.IsDeleted),
		Orgs:      orgs,
	})
}
