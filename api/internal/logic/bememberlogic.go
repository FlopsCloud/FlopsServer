package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type BeMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBeMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BeMemberLogic {
	return &BeMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BeMemberLogic) BeMember() (resp *types.Response, err error) {

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(uid))
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "加入会员失败",
			Info:    err.Error(),
		}, nil
	}

	if user.SysRole == "member" || user.SysRole == "superadmin" || user.SysRole == "admin" {
		return &types.Response{
			Code:    response.SuccessCode,
			Message: "你已经是成员",
		}, nil
	}

	user.SysRole = "member"
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "加入会员失败",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "成功加入会员",
	}, nil
}
