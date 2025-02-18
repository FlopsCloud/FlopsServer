package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateObjectLogic {
	return &UpdateObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateObjectLogic) UpdateObject(req *types.UpdateObjectRequest) (resp *types.ObjectResponse, err error) {
	// Check if object exists
	existingObject, err := l.svcCtx.ObjectsModel.FindOne(l.ctx, req.ObjId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.ObjectResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Object not found",
				},
			}, nil
		}
		return &types.ObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get object",
				Info:    err.Error(),
			},
		}, nil
	}

	// Update fields
	existingObject.Key = req.Key
	existingObject.Path = req.Path

	err = l.svcCtx.ObjectsModel.Update(l.ctx, existingObject)
	if err != nil {
		return &types.ObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to update object",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get updated object
	updatedObject, err := l.svcCtx.ObjectsModel.FindOne(l.ctx, req.ObjId)
	if err != nil {
		return &types.ObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Object updated but failed to retrieve",
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
			ObjId:     updatedObject.ObjId,
			BucketId:  updatedObject.BucketId,
			UserId:    updatedObject.UserId,
			Key:       updatedObject.Key,
			Path:      updatedObject.Path,
			CreatedAt: updatedObject.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: updatedObject.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
