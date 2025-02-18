package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TicketRepliesModel = (*customTicketRepliesModel)(nil)

type (
	// TicketRepliesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTicketRepliesModel.
	TicketRepliesModel interface {
		ticketRepliesModel
		withSession(session sqlx.Session) TicketRepliesModel
		FindByTicketId(ctx context.Context, ticketId uint64) ([]*TicketReplies, error)
		FindViewByTicketId(ctx context.Context, ticketId uint64) ([]*ViewTicketReplies, error)
	}

	customTicketRepliesModel struct {
		*defaultTicketRepliesModel
	}
)

type ViewTicketReplies struct {
	TicketReplies
	Nickname string `db:"nickname"`
	HeadUrl  string `db:"head_url"`
}

// NewTicketRepliesModel returns a model for the database table.
func NewTicketRepliesModel(conn sqlx.SqlConn) TicketRepliesModel {
	return &customTicketRepliesModel{
		defaultTicketRepliesModel: newTicketRepliesModel(conn),
	}
}

func (m *customTicketRepliesModel) withSession(session sqlx.Session) TicketRepliesModel {
	return NewTicketRepliesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTicketRepliesModel) FindByTicketId(ctx context.Context, ticketId uint64) ([]*TicketReplies, error) {
	query := fmt.Sprintf("select %s from %s where ticket_id = ? order by created_at asc", ticketRepliesRows, m.table)
	var replies []*TicketReplies
	err := m.conn.QueryRowsCtx(ctx, &replies, query, ticketId)
	return replies, err
}

func (m *customTicketRepliesModel) FindViewByTicketId(ctx context.Context, ticketId uint64) ([]*ViewTicketReplies, error) {
	query := fmt.Sprintf("SELECT a.*,b.nickname,b.head_url FROM ticket_replies a left outer join users b on a.user_id = b.user_id where ticket_id = ? order by created_at asc")
	var replies []*ViewTicketReplies
	err := m.conn.QueryRowsCtx(ctx, &replies, query, ticketId)
	return replies, err
}
