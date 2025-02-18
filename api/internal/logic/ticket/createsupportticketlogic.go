package ticket

import (
	"context"
	"encoding/json"
	"errors"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"fca/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSupportTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSupportTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSupportTicketLogic {
	return &CreateSupportTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSupportTicketLogic) CreateSupportTicket(req *types.CreateSupportTicketRequest) (resp *types.SupportTicketResponse, err error) {
	// 提取用户ID

	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		resp.Response = types.Response{
			Code:    response.ServerErrorCode,
			Message: err.Error(),
		}
		return resp, nil
	}
	userId := uint64(uid)

	// Validate priority
	if req.Priority == "" {
		req.Priority = "medium"
	}

	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
		"urgent": true,
	}
	if !validPriorities[req.Priority] {
		return nil, errors.New("invalid priority level")
	}

	// Create ticket
	ticket := &model.SupportTickets{
		UserId:      userId,
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
		Priority:    req.Priority,
		Images:      req.Images,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := l.svcCtx.SupportTicketsModel.Insert(l.ctx, ticket)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	ticket.TicketId = uint64(id)

	return &types.SupportTicketResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: convertToTicketType(ticket),
	}, nil
}
