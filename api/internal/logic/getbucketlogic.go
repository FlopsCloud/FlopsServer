package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBucketLogic {
	return &GetBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBucketLogic) GetBucket(req *types.GetBucketRequest) (resp *types.BucketResponse, err error) {
	bucket, err := l.svcCtx.BucketsModel.FindOne(l.ctx, req.BucketId)
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

	return &types.BucketResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.Bucket{
			BucketId:   bucket.BucketId,
			BucketName: bucket.BucketName,
			UserId:     bucket.UserId,
			Region:     bucket.Region,
			IsDeleted:  bucket.IsDeleted,
			DeletedAt:  bucket.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:  bucket.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  bucket.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
