package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BucketsModel = (*customBucketsModel)(nil)

type (
	// BucketsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBucketsModel.
	BucketsModel interface {
		bucketsModel
		withSession(session sqlx.Session) BucketsModel
		Count(ctx context.Context, region string) (uint64, error)
		FindMany(ctx context.Context, region string, offset, limit uint64) ([]*Buckets, error)
		FindByName(ctx context.Context, bucketName string) (*Buckets, error)
	}

	customBucketsModel struct {
		*defaultBucketsModel
	}
)

// NewBucketsModel returns a model for the database table.
func NewBucketsModel(conn sqlx.SqlConn) BucketsModel {
	return &customBucketsModel{
		defaultBucketsModel: newBucketsModel(conn),
	}
}

func (m *customBucketsModel) withSession(session sqlx.Session) BucketsModel {
	return NewBucketsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBucketsModel) FindByName(ctx context.Context, bucketName string) (*Buckets, error) {
	query := fmt.Sprintf("select %s from %s where `bucket_name` = ? and `is_deleted` = 0 limit 1", bucketsRows, m.table)
	var resp Buckets
	err := m.conn.QueryRowCtx(ctx, &resp, query, bucketName)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customBucketsModel) Count(ctx context.Context, region string) (uint64, error) {
	var count uint64
	query := fmt.Sprintf("select count(*) from %s where `is_deleted` = 0", m.table)

	if region != "" {
		query += " and `region` = ?"
		err := m.conn.QueryRowCtx(ctx, &count, query, region)
		return count, err
	}
	err := m.conn.QueryRowCtx(ctx, &count, query)
	return count, err
}

func (m *customBucketsModel) FindMany(ctx context.Context, region string, offset, limit uint64) ([]*Buckets, error) {
	query := fmt.Sprintf("select %s from %s where `is_deleted` = 0", bucketsRows, m.table)

	if region != "" {
		query += " and `region` = ?"
	}
	query += " order by `bucket_id` desc limit ?,?"

	var resp []*Buckets
	if region != "" {
		err := m.conn.QueryRowsCtx(ctx, &resp, query, region, offset, limit)
		return resp, err
	}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, offset, limit)
	return resp, err
}
