package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagListLogic {
	return &TagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagListLogic) TagList() (resp *types.TagsResp, err error) {
	tags, err := l.svcCtx.TagsModel.FindAll(l.ctx)
	if err != nil {
		return &types.TagsResp{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to fetch tags: " + err.Error(),
			},
		}, nil
	}

	var tagList []types.Tag
	for _, t := range tags {
		tagList = append(tagList, types.Tag{
			TagID:   t.TagId,
			TagName: t.TagName,
			// CreatedAt: t.CreatedAt.Unix(),
			// UpdatedAt: t.UpdatedAt.Unix(),
		})
	}

	return &types.TagsResp{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: tagList,
	}, nil
}
