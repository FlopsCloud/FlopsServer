package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagAddLogic {
	return &TagAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagAddLogic) TagAdd(req *types.TagAddRequest) (resp *types.TagsResp, err error) {
	// Check if tag name already exists
	exists, err := l.svcCtx.TagsModel.FindByName(l.ctx, req.TagName)
	if err != nil && err != model.ErrNotFound {
		return &types.TagsResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: err.Error(),
			},
		}, nil
	}
	if exists != nil {
		return &types.TagsResp{
			Response: types.Response{
				Code:    response.ParameterErrorCode,
				Message: "Tag name already exists",
			},
		}, nil
	}

	// Create new tag
	_, err = l.svcCtx.TagsModel.Insert(l.ctx, &model.Tags{
		TagName:   req.TagName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return &types.TagsResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: err.Error(),
			},
		}, nil
	}

	return &types.TagsResp{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Tag added successfully",
		},
	}, nil
}
