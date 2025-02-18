package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteObjectLogic {
	return &DeleteObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteObjectLogic) DeleteObject(req *types.DeleteObjectRequest) (resp *types.Response, err error) {
	// First find the bucket by name
	bucket, err := l.svcCtx.BucketsModel.FindByName(l.ctx, req.BucketName)
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

	// Find object by bucket ID and key
	object, err := l.svcCtx.ObjectsModel.FindByBucketAndKey(l.ctx, bucket.BucketId, req.ObjectKey)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: "Object not found",
			}, nil
		}
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to get object",
			Info:    err.Error(),
		}, nil
	}

	// Delete object
	err = l.svcCtx.ObjectsModel.Delete(l.ctx, object.ObjId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "Failed to delete object",
			Info:    err.Error(),
		}, nil
	}

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "Success",
	}, nil
}
