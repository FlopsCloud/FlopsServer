package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateImageLogic {
	return &UpdateImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateImageLogic) UpdateImage(req *types.UpdateImageRequest) (resp *types.ImageResponse, err error) {
	// Check if image exists
	existingImage, err := l.svcCtx.ImagesModel.FindOne(l.ctx, req.ImageId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.ImageResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Image not found",
				},
			}, nil
		}
		return &types.ImageResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get image",
				Info:    err.Error(),
			},
		}, nil
	}

	// Update image
	existingImage.ImageName = req.ImageName
	err = l.svcCtx.ImagesModel.Update(l.ctx, existingImage)
	if err != nil {
		return &types.ImageResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to update image",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get updated image
	updatedImage, err := l.svcCtx.ImagesModel.FindOne(l.ctx, req.ImageId)
	if err != nil {
		return &types.ImageResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Image updated but failed to retrieve",
				Info:    err.Error(),
			},
		}, nil
	}

	return &types.ImageResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.Image{
			ImageId:   updatedImage.ImageId,
			ImageName: updatedImage.ImageName,
			CreatedAt: updatedImage.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: updatedImage.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
