package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBucketLogic {
	return &UpdateBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBucketLogic) UpdateBucket(req *types.UpdateBucketRequest) (resp *types.BucketResponse, err error) {
	// Check if bucket exists
	existingBucket, err := l.svcCtx.BucketsModel.FindOne(l.ctx, req.BucketId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BucketResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Bucket not found",
				},
			}, nil
		}
		return &types.BucketResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get bucket",
				Info:    err.Error(),
			},
		}, nil
	}

	// Update fields
	existingBucket.BucketName = req.BucketName
	existingBucket.Region = req.Region

	err = l.svcCtx.BucketsModel.Update(l.ctx, existingBucket)
	if err != nil {
		return &types.BucketResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to update bucket",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get updated bucket
	updatedBucket, err := l.svcCtx.BucketsModel.FindOne(l.ctx, req.BucketId)
	if err != nil {
		return &types.BucketResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Bucket updated but failed to retrieve",
				Info:    err.Error(),
			},
		}, nil
	}

	return &types.BucketResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.Bucket{
			BucketId:   updatedBucket.BucketId,
			BucketName: updatedBucket.BucketName,
			UserId:     updatedBucket.UserId,
			Region:     updatedBucket.Region,
			IsDeleted:  updatedBucket.IsDeleted,
			DeletedAt:  updatedBucket.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:  updatedBucket.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  updatedBucket.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
