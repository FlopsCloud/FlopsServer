package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ InvitationModel = (*customInvitationModel)(nil)

type (
	// InvitationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInvitationModel.
	InvitationModel interface {
		invitationModel
		withSession(session sqlx.Session) InvitationModel
	}

	customInvitationModel struct {
		*defaultInvitationModel
	}
)

// NewInvitationModel returns a model for the database table.
func NewInvitationModel(conn sqlx.SqlConn) InvitationModel {
	return &customInvitationModel{
		defaultInvitationModel: newInvitationModel(conn),
	}
}

func (m *customInvitationModel) withSession(session sqlx.Session) InvitationModel {
	return NewInvitationModel(sqlx.NewSqlConnFromSession(session))
}
