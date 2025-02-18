package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ObjectsTreeModel = (*customObjectsTreeModel)(nil)

type (
	// ObjectsTreeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customObjectsTreeModel.
	ObjectsTreeModel interface {
		objectsTreeModel
		WithSession(session sqlx.Session) ObjectsTreeModel
		FindByBucketAndFullName(ctx context.Context, bucketId uint64, fullTreeName string) (*ObjectsTree, error)
		FindAllByBucketId(ctx context.Context, bucketId uint64) (*[]ObjectsTree, error)
	}

	customObjectsTreeModel struct {
		*defaultObjectsTreeModel
	}
)

// NewObjectsTreeModel returns a model for the database table.
func NewObjectsTreeModel(conn sqlx.SqlConn) ObjectsTreeModel {
	return &customObjectsTreeModel{
		defaultObjectsTreeModel: newObjectsTreeModel(conn),
	}
}

func (m *customObjectsTreeModel) WithSession(session sqlx.Session) ObjectsTreeModel {
	return NewObjectsTreeModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customObjectsTreeModel) FindByBucketAndFullName(ctx context.Context, bucketId uint64, fullTreeName string) (*ObjectsTree, error) {
	query := fmt.Sprintf("select %s from %s where `bucket_id` = ? and `full_tree_name` = ? limit 1", objectsTreeRows, m.table)
	var resp ObjectsTree
	err := m.conn.QueryRowCtx(ctx, &resp, query, bucketId, fullTreeName)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customObjectsTreeModel) FindAllByBucketId(ctx context.Context, bucketId uint64) (*[]ObjectsTree, error) {
	query := fmt.Sprintf("select %s from %s where `bucket_id` = ? order by `parant_tree_id`", objectsTreeRows, m.table)
	var resp []ObjectsTree
	err := m.conn.QueryRowsCtx(ctx, &resp, query, bucketId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return &resp, nil
	default:
		return nil, err
	}
}
