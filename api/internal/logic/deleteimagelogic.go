package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteImageLogic {
	return &DeleteImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteImageLogic) DeleteImage(req *types.DeleteImageRequest) (resp *types.Response, err error) {
	// Check if image exists
	_, err = l.svcCtx.ImagesModel.FindOne(l.ctx, req.ImageId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: "Image not found",
			}, nil
		}
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to get image",
			Info:    err.Error(),
		}, nil
	}

	// Delete image
	err = l.svcCtx.ImagesModel.Delete(l.ctx, req.ImageId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to delete image",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}, nil
}
