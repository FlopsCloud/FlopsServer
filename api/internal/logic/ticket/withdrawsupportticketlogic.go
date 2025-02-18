package ticket

import (
	"context"
	"fmt"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawSupportTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// user and admin use same api，images is comma separated
func NewWithdrawSupportTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawSupportTicketLogic {
	return &WithdrawSupportTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func NewResp(code int64, message, info string) *types.Response {
	return &types.Response{

		Code:    code,
		Message: message,
		Info:    info,
	}
}

func (l *WithdrawSupportTicketLogic) WithdrawSupportTicket(req *types.WithdrawSupportTicketRequest) (resp *types.Response, err error) {
	// Get the reply
	reply, err := l.svcCtx.TicketRepliesModel.FindOne(l.ctx, req.ReplyId)
	if err != nil {
		if err == model.ErrNotFound {
			return NewResp(response.ServerErrorCode, "回复不存在", ""), nil
		}
		return NewResp(response.ServerErrorCode, "回复不存在", err.Error()), nil
	}

	// Verify the reply belongs to the specified ticket
	if reply.TicketId != req.TicketId {
		return NewResp(response.ServerErrorCode, "回复不属于该工单", ""), nil
	}

	// Check if the reply is within 2 minutes

	if time.Since(reply.CreatedAt) > 2*time.Minute {
		return NewResp(response.InvalidRequestParamCode, fmt.Sprintf("不能在2分钟内撤回回复. 当前回复时间: %.2f 分钟",
			time.Since(reply.CreatedAt).Minutes()), ""), nil
	}

	// Delete the reply
	err = l.svcCtx.TicketRepliesModel.Delete(l.ctx, req.ReplyId)
	if err != nil {
		return NewResp(response.ServerErrorCode, "撤回回复失败", err.Error()), nil
	}

	return NewResp(response.SuccessCode, "撤回回复成功", ""), nil
}
