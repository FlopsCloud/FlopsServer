package logic

import (
	"context"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddServerTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddServerTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddServerTagLogic {
	return &AddServerTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddServerTagLogic) AddServerTag(req *types.AddServerTagRequest) (resp *types.Response, err error) {
	serverTag := &model.ServerTags{
		ServerId:  req.ServerId,
		TagId:     uint64(req.TagID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = l.svcCtx.ServerTagsModel.Insert(l.ctx, serverTag)
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: "Fail",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    0,
		Message: "Success",
	}, nil
}
