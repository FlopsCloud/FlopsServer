package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		WithSession(session sqlx.Session) UsersModel
		RestPass(ctx context.Context, data *Users) error
		FindOneByPhone(ctx context.Context, phone string) (*Users, error)
		FindUsers(ctx context.Context, username, email, phone string, page, pageSize uint64) ([]*Users, uint64, error)
		UpdateInfo(ctx context.Context, newData *Users) error
		FindOneByX(ctx context.Context, AccX string) (*Users, error)
		FindOneByGoogle(ctx context.Context, AccGoogle string) (*Users, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) WithSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUsersModel) FindOneByX(ctx context.Context, AccX string) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `acc_x` = ? and is_deleted=0 limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, AccX)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindOneByGoogle(ctx context.Context, AccGoogle string) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `acc_google` = ? and is_deleted=0 limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, AccGoogle)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindOneByEmail(ctx context.Context, email string) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `email` = ? and is_deleted=0 limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) UpdateInfo(ctx context.Context, newData *Users) error {

	usersRowsWithPlaceHolder2 := strings.Join([]string{"`username`", "`nickname`", "`email`", "`phone`"}, "=?,") + "=?"

	logx.Infof("update %s set %s where `user_id` = ?", m.table, usersRowsWithPlaceHolder2)

	query := fmt.Sprintf("update %s set %s where `user_id` = ?", m.table, usersRowsWithPlaceHolder2)
	_, err := m.conn.ExecCtx(ctx, query, newData.Username, newData.Nickname, newData.Email, newData.Phone, newData.UserId)
	return err
}

func (m *customUsersModel) FindOneByPhone(ctx context.Context, phone string) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `phone` = ? and is_deleted=0 limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) Delete(ctx context.Context, userId uint64) error {
	query := fmt.Sprintf("update %s set is_deleted=1 where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

func (m *customUsersModel) FindOne(ctx context.Context, userId uint64) (*Users, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and is_deleted=0 limit 1", usersRows, m.table)
	var resp Users
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) RestPass(ctx context.Context, newData *Users) error {
	query := fmt.Sprintf("update %s set password_hash=? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, newData.PasswordHash, newData.UserId)
	return err
}

func (m *customUsersModel) FindUsers(ctx context.Context, username, email, phone string, page, pageSize uint64) ([]*Users, uint64, error) {
	//whereClause := "WHERE is_deleted = 0"
	whereClause := "WHERE 1 = 1"
	params := []interface{}{}

	if username != "" {
		whereClause += " AND username LIKE ?"
		params = append(params, fmt.Sprintf("%%%s%%", username))
	}
	if email != "" {
		whereClause += " AND email LIKE ?"
		params = append(params, fmt.Sprintf("%%%s%%", email))
	}
	if phone != "" {
		whereClause += " AND phone LIKE ?"
		params = append(params, fmt.Sprintf("%%%s%%", phone))
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereClause)
	var total uint64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT %s FROM %s %s LIMIT ? OFFSET ?", usersRows, m.table, whereClause)
	params = append(params, pageSize, offset)

	var users []*Users
	err = m.conn.QueryRowsCtx(ctx, &users, query, params...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
