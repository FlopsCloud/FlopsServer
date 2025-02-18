package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserListLogic {
	return &AdminUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserListLogic) AdminUserList(req *types.AdminUserListRequest) (resp *types.AdminUserListResponse, err error) {
	resp = &types.AdminUserListResponse{
		Response: types.Response{
			Code:    0,
			Message: "Success",
		},
	}

	// 管理员权限验证 start
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	email, _ := l.ctx.Value("email").(string)

	logx.Info("JWT uid=", uid, " Name=", email)
	if uid != 1 {
		resp.Code = 1
		resp.Message = "Permission denied"
		return resp, nil
	}
	// 管理员权限验证 end

	users, total, err := l.svcCtx.UserModel.FindUsers(l.ctx, req.Username, req.Email, req.Phone, req.Page, req.PageSize)
	if err != nil {
		resp.Code = 1
		resp.Message = "Failed to fetch users: " + err.Error()
		return resp, err
	}

	userList := make([]types.UserResponseData, len(users))
	for i, user := range users {
		userList[i] = types.UserResponseData{
			UserId:    user.UserId,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Phone:     user.Phone,
			Email:     user.Email,
			Balance:   user.Balance,
			IsDeleted: int64(user.IsDeleted),

			//TODO : add org
			// Note: Orgs field is not populated here. If needed, you should fetch and populate it separately.
		}
	}

	resp.Data = types.AdminUserListResponseData{
		Users: userList,
		Total: total,
	}

	return resp, nil
}
