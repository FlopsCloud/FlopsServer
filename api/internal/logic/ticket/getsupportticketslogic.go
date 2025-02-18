package ticket

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSupportTicketsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSupportTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSupportTicketsLogic {
	return &GetSupportTicketsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSupportTicketsLogic) GetSupportTickets(req *types.GetSupportTicketsRequest) (resp *types.GetSupportTicketsResponse, err error) {
	// Get the ticket
	ticket, err := l.svcCtx.SupportTicketsModel.FindOne(l.ctx, req.TicketId)
	if err != nil {
		return &types.GetSupportTicketsResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get ticket",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get all replies for this ticket
	replies, err := l.svcCtx.TicketRepliesModel.FindViewByTicketId(l.ctx, req.TicketId)
	if err != nil {
		return &types.GetSupportTicketsResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get ticket replies",
				Info:    err.Error(),
			},
		}, nil
	}

	// Convert replies to response type
	replyList := make([]types.ReplyTicket, len(replies))
	for i, r := range replies {
		replyList[i] = types.ReplyTicket{
			ReplyId:   r.ReplyId,
			TicketId:  r.TicketId,
			UserId:    r.UserId,
			Content:   r.Content,
			Nickname:  r.Nickname,
			HeadUrl:   r.HeadUrl,
			Images:    r.Images,
			CreatedAt: uint64(r.CreatedAt.Unix()),
		}
	}

	return &types.GetSupportTicketsResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: types.GetSupportTicketsResponseData{
			Ticket:      convertToTicketType(ticket),
			ReplyTicket: replyList,
		},
	}, nil
}
