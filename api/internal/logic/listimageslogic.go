package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListImagesLogic {
	return &ListImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListImagesLogic) ListImages(req *types.ListImagesRequest) (resp *types.ListImagesResponse, err error) {
	// Set default page size if not provided
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Get total count
	total, err := l.svcCtx.ImagesModel.Count(l.ctx, req.ImageName)
	if err != nil {
		return &types.ListImagesResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get total count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get images with pagination
	images, err := l.svcCtx.ImagesModel.FindMany(l.ctx, req.ImageName, offset, req.PageSize)
	if err != nil {
		return &types.ListImagesResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get images",
				Info:    err.Error(),
			},
		}, nil
	}

	// Convert to response type
	var imageList []types.Image
	for _, img := range images {
		imageList = append(imageList, types.Image{
			ImageId:   img.ImageId,
			ImageName: img.ImageName,
			CreatedAt: img.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: img.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ListImagesResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ListImagesResponseData{
			Images: imageList,
			Total:  total,
		},
	}, nil
}
