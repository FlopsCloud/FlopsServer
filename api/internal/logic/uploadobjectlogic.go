package logic

import (
	"context"
	"encoding/json"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadObjectLogic {
	return &UploadObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadObjectLogic) UploadObject(req *types.UploadObjectRequest) (resp *types.ObjectResponse, err error) {
	// First find the bucket by name
	bucket, err := l.svcCtx.BucketsModel.FindByName(l.ctx, req.BucketName)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.ObjectResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Bucket not found",
				},
			}, nil
		}
		return &types.ObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get bucket",
				Info:    err.Error(),
			},
		}, nil
	}

	// Check if object with same key already exists
	existingObject, err := l.svcCtx.ObjectsModel.FindByBucketAndKey(l.ctx, bucket.BucketId, req.ObjectKey)
	if err != nil && err != model.ErrNotFound {
		return &types.ObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to check existing object",
				Info:    err.Error(),
			},
		}, nil
	}
	if existingObject != nil {
		return &types.ObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Object with this key already exists",
			},
		}, nil
	}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	// Create new object
	object := &model.Objects{
		BucketId: bucket.BucketId,
		UserId:   uint64(uid),
		Key:      req.ObjectKey,
		Path:     "/data/" + bucket.BucketName + "/" + req.ObjectKey, // Simple path construction
	}

	result, err := l.svcCtx.ObjectsModel.Insert(l.ctx, object)
	if err != nil {
		return &types.ObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to create object",
				Info:    err.Error(),
			},
		}, nil
	}

	objId, _ := result.LastInsertId()
	createdObject, err := l.svcCtx.ObjectsModel.FindOne(l.ctx, uint64(objId))
	if err != nil {
		return &types.ObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Object created but failed to retrieve",
				Info:    err.Error(),
			},
		}, nil
	}

	// TODO: Actually store the content somewhere (file system, object storage, etc.)
	// For now, we just store the metadata in the database

	return &types.ObjectResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.Object{
			ObjId:     createdObject.ObjId,
			BucketId:  createdObject.BucketId,
			UserId:    createdObject.UserId,
			Key:       createdObject.Key,
			Path:      createdObject.Path,
			CreatedAt: createdObject.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: createdObject.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
