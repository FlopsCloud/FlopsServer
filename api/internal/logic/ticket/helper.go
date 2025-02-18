package ticket

import (
	"fca/api/internal/types"
	"fca/model"
)

func convertToTicketType(ticket *model.SupportTickets) types.SupportTicket {
	return types.SupportTicket{
		TicketId:    ticket.TicketId,
		UserId:      ticket.UserId,
		Title:       ticket.Title,
		Description: ticket.Description,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		Images:      ticket.Images,
		CreatedAt:   ticket.CreatedAt.Unix(),
		UpdatedAt:   ticket.UpdatedAt.Unix(),
	}
}
