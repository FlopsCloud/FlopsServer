package logic

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/cryptx"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPassLogic {
	return &ResetPassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPassLogic) ResetPass(req *types.RestPassRequest) (resp response.Response, err error) {
	// // Extract userId from context

	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return response.Fail(response.UnauthorizedCode, "Invalid user ID"), nil

	}
	userId := uint64(uid)

	logx.Infof("ResetPass userId: %d", userId)

	// Find the user by userId
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == model.ErrNotFound {
			return response.Fail(response.UserNotExistCode, "User not found"), nil
		}
		return response.Error("Failed to find user"), err
	}

	// Verify that the email in the request matches the user's email
	// if user.Email != req.Email {
	// 	return response.Fail(response.UnauthorizedCode, "Email does not match authenticated user"), nil
	// }

	// Verify old password
	// oldPasswordHash := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, req.OldPassword)
	// if oldPasswordHash != user.PasswordHash {
	// 	return response.Fail(response.UnauthorizedCode, "Invalid old password"), nil
	// }

	// Hash the new password
	newPasswordHash := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, req.NewPassword)

	// Update the user's password
	err = l.svcCtx.UserModel.RestPass(l.ctx, &model.Users{
		UserId:       user.UserId,
		PasswordHash: newPasswordHash,
	})
	if err != nil {
		return response.Error("Failed to update password"), err
	}

	// Return success response
	resp.Code = 200
	resp.Message = "Password reset successfully"
	return resp, nil

}
