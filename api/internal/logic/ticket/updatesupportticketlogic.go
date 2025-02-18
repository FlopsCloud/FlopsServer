package ticket

import (
	"context"
	"errors"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSupportTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSupportTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSupportTicketLogic {
	return &UpdateSupportTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSupportTicketLogic) UpdateSupportTicket(req *types.UpdateSupportTicketRequest) (resp *types.SupportTicketResponse, err error) {
	// l.Logger.Info("UpdateSupportTicketLogic: UpdateSupportTicket", "req", req)

	// Get existing ticket
	ticket, err := l.svcCtx.SupportTicketsModel.FindOne(l.ctx, req.TicketId)
	if err != nil {
		return nil, err
	}

	// Validate status if provided
	if req.Status != "" {
		validStatuses := map[string]bool{
			"open":        true,
			"in-progress": true,
			"resolved":    true,
			"closed":      true,
		}
		if !validStatuses[req.Status] {
			return nil, errors.New("invalid status")
		}
		ticket.Status = req.Status
	}

	// Update description if provided
	if req.Description != "" {
		ticket.Description = req.Description
	}

	ticket.UpdatedAt = time.Now()

	err = l.svcCtx.SupportTicketsModel.Update(l.ctx, ticket)
	if err != nil {
		return nil, err
	}

	return &types.SupportTicketResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: convertToTicketType(ticket),
	}, nil
}
