package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteResourceLogic {
	return &DeleteResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteResourceLogic) DeleteResource(resourceId uint64) (resp *types.Response, err error) {
	// Check if resource exists
	// l.Logger.Info("in DeleteResourceLogic")
	_, err = l.svcCtx.ResourcesModel.FindOne(l.ctx, resourceId)
	if err != nil {
		return &types.Response{
			Code:    404,
			Message: "Not found",
			Info:    err.Error(),
		}, nil
	}

	// Soft delete the resource
	err = l.svcCtx.ResourcesModel.Delete(l.ctx, resourceId)
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: "Delete Error",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "success",
	}, nil
}
