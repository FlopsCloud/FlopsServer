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

type CreateObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateObjectLogic {
	return &CreateObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateObjectLogic) CreateObject(req *types.CreateObjectRequest) (resp *types.ObjectResponse, err error) {
	// Check if bucket exists
	_, err = l.svcCtx.BucketsModel.FindOne(l.ctx, req.BucketId)
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
				Message: "Failed to check bucket",
				Info:    err.Error(),
			},
		}, nil
	}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	object := &model.Objects{
		BucketId: req.BucketId,
		Key:      req.Key,
		Path:     req.Path,
		UserId:   uint64(uid),
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
