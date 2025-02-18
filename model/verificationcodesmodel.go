package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VerificationCodesModel = (*customVerificationCodesModel)(nil)

type (
	// VerificationCodesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVerificationCodesModel.
	VerificationCodesModel interface {
		verificationCodesModel
		withSession(session sqlx.Session) VerificationCodesModel
		FindOneByPhone(ctx context.Context, phone string) (*VerificationCodes, error)
		FindOneByEmail(ctx context.Context, email string) (*VerificationCodes, error)
	}

	customVerificationCodesModel struct {
		*defaultVerificationCodesModel
	}
)

// NewVerificationCodesModel returns a model for the database table.
func NewVerificationCodesModel(conn sqlx.SqlConn) VerificationCodesModel {
	return &customVerificationCodesModel{
		defaultVerificationCodesModel: newVerificationCodesModel(conn),
	}
}

func (m *customVerificationCodesModel) withSession(session sqlx.Session) VerificationCodesModel {
	return NewVerificationCodesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customVerificationCodesModel) FindOneByPhone(ctx context.Context, phone string) (*VerificationCodes, error) {
	query := fmt.Sprintf("select %s from %s where `target_value` = ? and target_type=? order by id desc limit 1", verificationCodesRows, m.table)
	var resp VerificationCodes
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone, "phone")
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customVerificationCodesModel) FindOneByEmail(ctx context.Context, email string) (*VerificationCodes, error) {
	query := fmt.Sprintf("select %s from %s where `target_value` = ?  and target_type=? order by id desc limit 1", verificationCodesRows, m.table)
	var resp VerificationCodes
	err := m.conn.QueryRowCtx(ctx, &resp, query, email, "email")
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
