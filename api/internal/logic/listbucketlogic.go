package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBucketLogic {
	return &ListBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBucketLogic) ListBucket(req *types.ListBucketRequest) (resp *types.ListBucketResponse, err error) {
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
	total, err := l.svcCtx.BucketsModel.Count(l.ctx, req.Region)
	if err != nil {
		return &types.ListBucketResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get total count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get buckets with pagination
	buckets, err := l.svcCtx.BucketsModel.FindMany(l.ctx, req.Region, offset, req.PageSize)
	if err != nil {
		return &types.ListBucketResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get buckets",
				Info:    err.Error(),
			},
		}, nil
	}

	// Convert to response type
	var bucketList []types.Bucket
	for _, b := range buckets {
		bucketList = append(bucketList, types.Bucket{
			BucketId:   b.BucketId,
			BucketName: b.BucketName,
			UserId:     b.UserId,
			Region:     b.Region,
			IsDeleted:  b.IsDeleted,
			DeletedAt:  b.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:  b.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  b.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ListBucketResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ListBucketResponseData{
			Buckets: bucketList,
			Total:   total,
		},
	}, nil
}
