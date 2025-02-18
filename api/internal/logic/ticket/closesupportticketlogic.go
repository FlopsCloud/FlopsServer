package ticket

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloseSupportTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// user close ticket
func NewCloseSupportTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseSupportTicketLogic {
	return &CloseSupportTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloseSupportTicketLogic) CloseSupportTicket(req *types.CloseSupportTicketRequest) (resp *types.Response, err error) {

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	// role, _ := l.ctx.Value("role").(string)

	ticket, err := l.svcCtx.SupportTicketsModel.FindOne(l.ctx, req.TicketId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "找不到工单",
			Info:    err.Error(),
		}, nil
	}
	if ticket.UserId != uint64(uid) {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "不能关闭他人的工单",
		}, nil
	}

	ticket.Status = "closed"
	err = l.svcCtx.SupportTicketsModel.Update(l.ctx, ticket)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "关闭工单失败",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "关闭工单成功",
	}, nil
}
