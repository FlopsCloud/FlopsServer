package logic

import (
	"context"
	"database/sql"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBucketLogic {
	return &DeleteBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBucketLogic) DeleteBucket(req *types.DeleteBucketRequest) (resp *types.Response, err error) {
	// Check if bucket exists
	existingBucket, err := l.svcCtx.BucketsModel.FindOne(l.ctx, req.BucketId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: "Bucket not found",
			}, nil
		}
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to get bucket",
			Info:    err.Error(),
		}, nil
	}

	// Soft delete by updating fields
	existingBucket.IsDeleted = 1
	existingBucket.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	err = l.svcCtx.BucketsModel.Update(l.ctx, existingBucket)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to delete bucket",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}, nil
}
