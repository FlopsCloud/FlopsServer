package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagDeleteLogic {
	return &TagDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagDeleteLogic) TagDelete(req *types.TagDeleteRequest) (resp *types.TagsResp, err error) {
	// Check if tag exists
	_, err = l.svcCtx.TagsModel.FindOne(l.ctx, req.TagId)
	if err == model.ErrNotFound {
		return &types.TagsResp{
			Response: types.Response{
				Code:    response.NotFoundCode,
				Message: "Tag not found",
			},
		}, nil
	}
	if err != nil {
		return &types.TagsResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: err.Error(),
			},
		}, nil
	}

	// Delete tag
	err = l.svcCtx.TagsModel.Delete(l.ctx, req.TagId)
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
			Message: "Tag deleted successfully",
		},
	}, nil
}
