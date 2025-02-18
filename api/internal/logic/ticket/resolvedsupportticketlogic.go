package ticket

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResolvedSupportTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// user resolved ticket
func NewResolvedSupportTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolvedSupportTicketLogic {
	return &ResolvedSupportTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResolvedSupportTicketLogic) ResolvedSupportTicket(req *types.CloseSupportTicketRequest) (resp *types.Response, err error) {

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
			Message: "不能他人的工单为已解决",
		}, nil
	}

	ticket.Status = "resolved"
	err = l.svcCtx.SupportTicketsModel.Update(l.ctx, ticket)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "解决工单失败",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "解决工单成功",
	}, nil
}
