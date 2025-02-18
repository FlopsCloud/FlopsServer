package model

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SupportTicketsModel = (*customSupportTicketsModel)(nil)

type (
	// SupportTicketsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSupportTicketsModel.
	SupportTicketsModel interface {
		supportTicketsModel
		withSession(session sqlx.Session) SupportTicketsModel
		FindList(whereBuilder string, pageSize, offset int64, args ...interface{}) ([]*SupportTickets, error)
		Count(whereBuilder string, args ...interface{}) (int64, error)
	}

	customSupportTicketsModel struct {
		*defaultSupportTicketsModel
	}
)

// NewSupportTicketsModel returns a model for the database table.
func NewSupportTicketsModel(conn sqlx.SqlConn) SupportTicketsModel {
	return &customSupportTicketsModel{
		defaultSupportTicketsModel: newSupportTicketsModel(conn),
	}
}

func (m *customSupportTicketsModel) withSession(session sqlx.Session) SupportTicketsModel {
	return NewSupportTicketsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSupportTicketsModel) FindList(whereBuilder string, pageSize, offset int64, args ...interface{}) ([]*SupportTickets, error) {
	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY created_at DESC LIMIT ? OFFSET ?", m.table, whereBuilder)
	args = append(args, pageSize, offset)

	var tickets []*SupportTickets
	err := m.conn.QueryRows(&tickets, query, args...)
	return tickets, err
}

func (m *customSupportTicketsModel) Count(whereBuilder string, args ...interface{}) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereBuilder)
	var count int64
	err := m.conn.QueryRow(&count, query, args...)
	return count, err
}
