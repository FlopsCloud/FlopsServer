package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ImagesModel = (*customImagesModel)(nil)

type (
	// ImagesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customImagesModel.
	ImagesModel interface {
		imagesModel
		withSession(session sqlx.Session) ImagesModel
		Count(ctx context.Context, imageName string) (uint64, error)
		FindMany(ctx context.Context, imageName string, offset, limit uint64) ([]*Images, error)
	}

	customImagesModel struct {
		*defaultImagesModel
	}
)

// NewImagesModel returns a model for the database table.
func NewImagesModel(conn sqlx.SqlConn) ImagesModel {
	return &customImagesModel{
		defaultImagesModel: newImagesModel(conn),
	}
}

func (m *customImagesModel) withSession(session sqlx.Session) ImagesModel {
	return NewImagesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customImagesModel) Count(ctx context.Context, imageName string) (uint64, error) {
	var count uint64
	query := fmt.Sprintf("select count(*) from %s", m.table)
	if imageName != "" {
		query += " where `image_name` like ?"
		imageName = "%" + imageName + "%"
		err := m.conn.QueryRowCtx(ctx, &count, query, imageName)
		return count, err
	}
	err := m.conn.QueryRowCtx(ctx, &count, query)
	return count, err
}

func (m *customImagesModel) FindMany(ctx context.Context, imageName string, offset, limit uint64) ([]*Images, error) {
	query := fmt.Sprintf("select %s from %s", imagesRows, m.table)
	if imageName != "" {
		query += " where `image_name` like ?"
		imageName = "%" + imageName + "%"
	}
	query += " order by `image_id` desc limit ?,?"

	var resp []*Images
	if imageName != "" {
		err := m.conn.QueryRowsCtx(ctx, &resp, query, imageName, offset, limit)
		return resp, err
	}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, offset, limit)
	return resp, err
}
