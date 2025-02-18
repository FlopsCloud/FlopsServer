package ticket

import (
	"context"
	"encoding/json"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSupportTicketsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSupportTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSupportTicketsLogic {
	return &ListSupportTicketsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSupportTicketsLogic) ListSupportTickets(req *types.ListSupportTicketsRequest) (resp *types.ListSupportTicketsResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	// Build query conditions
	whereBuilder := strings.Builder{}
	args := []interface{}{}

	whereBuilder.WriteString(" WHERE 1=1")

	if req.Status != "" {
		whereBuilder.WriteString(" AND status = ?")
		args = append(args, req.Status)
	}

	if uid != 0 {
		whereBuilder.WriteString(" AND user_id = ?")
		args = append(args, uid)
	}

	// Add pagination
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}
	offset := (req.Page - 1) * pageSize

	// Convert uint64 to int64 for the model methods
	pageSizeInt64 := int64(pageSize)
	offsetInt64 := int64(offset)

	// Use the model methods
	total, err := l.svcCtx.SupportTicketsModel.Count(whereBuilder.String(), args...)
	if err != nil {
		return nil, err
	}

	tickets, err := l.svcCtx.SupportTicketsModel.FindList(whereBuilder.String(), pageSizeInt64, offsetInt64, args...)
	if err != nil {
		return nil, err
	}

	// Convert to response type
	ticketList := make([]types.SupportTicket, len(tickets))
	for i, t := range tickets {
		ticketList[i] = convertToTicketType(t)
	}

	return &types.ListSupportTicketsResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "success",
		},
		Data: types.ListSupportTicketsResponseData{
			Tickets: ticketList,
			Total:   uint64(total),
		},
	}, nil
}
