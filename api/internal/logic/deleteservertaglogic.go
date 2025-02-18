package logic

import (
	"context"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteServerTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteServerTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteServerTagLogic {
	return &DeleteServerTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteServerTagLogic) DeleteServerTag(req *types.DeleteServerTagRequest) (resp *types.Response, err error) {
	// Create a ServerTags object with the request data
	serverTag := &model.ServerTags{
		ServerId: req.ServerId,
		TagId:    uint64(req.TagID),
	}

	// Try to delete the record
	err = l.svcCtx.ServerTagsModel.DeleteByServerTag(l.ctx, serverTag)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    404,
				Message: fmt.Sprintf("Server-tag association not found for ServerId: %d and TagId: %d", req.ServerId, req.TagID),
			}, nil
		}
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "Success",
	}, nil
}
