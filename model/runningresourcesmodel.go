package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RunningResourcesModel = (*customRunningResourcesModel)(nil)

type (
	// RunningResourcesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRunningResourcesModel.
	RunningResourcesModel interface {
		runningResourcesModel
		withSession(session sqlx.Session) RunningResourcesModel
		FindAllRunning(ctx context.Context) ([]*RunningResources, error)
		FindAll(ctx context.Context, where string) ([]*RunningResources, error)
		Stop(ctx context.Context, instanceId uint64) error
	}

	customRunningResourcesModel struct {
		*defaultRunningResourcesModel
	}
)

// NewRunningResourcesModel returns a model for the database table.
func NewRunningResourcesModel(conn sqlx.SqlConn) RunningResourcesModel {
	return &customRunningResourcesModel{
		defaultRunningResourcesModel: newRunningResourcesModel(conn),
	}
}

func (m *customRunningResourcesModel) withSession(session sqlx.Session) RunningResourcesModel {
	return NewRunningResourcesModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customRunningResourcesModel) Stop(ctx context.Context, instanceId uint64) error {
	query := fmt.Sprintf("update %s set stat = 'stopped', stop_at = ? where instance_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), instanceId)
	return err
}

func (m *customRunningResourcesModel) FindAllRunning(ctx context.Context) ([]*RunningResources, error) {
	return m.FindAll(ctx, "stat = 'started'")
}

func (m *customRunningResourcesModel) FindAll(ctx context.Context, where string) ([]*RunningResources, error) {
	var resp []*RunningResources
	query := fmt.Sprintf("select %s from %s where %s", runningResourcesRows, m.table, where)
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	return resp, err
}
