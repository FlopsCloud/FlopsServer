package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TagsModel = (*customTagsModel)(nil)

type (
	// TagsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTagsModel.
	TagsModel interface {
		tagsModel
		withSession(session sqlx.Session) TagsModel
	}

	customTagsModel struct {
		*defaultTagsModel
	}
)

// NewTagsModel returns a model for the database table.
func NewTagsModel(conn sqlx.SqlConn) TagsModel {
	return &customTagsModel{
		defaultTagsModel: newTagsModel(conn),
	}
}

func (m *customTagsModel) withSession(session sqlx.Session) TagsModel {
	return NewTagsModel(sqlx.NewSqlConnFromSession(session))
}
