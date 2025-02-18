package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImageLogic {
	return &GetImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetImageLogic) GetImage(req *types.GetImageRequest) (resp *types.ImageResponse, err error) {
	image, err := l.svcCtx.ImagesModel.FindOne(l.ctx, req.ImageId)
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

	return &types.ImageResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.Image{
			ImageId:   image.ImageId,
			ImageName: image.ImageName,
			CreatedAt: image.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: image.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
