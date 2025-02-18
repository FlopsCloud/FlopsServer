package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetObjectLogic {
	return &GetObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetObjectLogic) GetObject(req *types.GetObjectRequest) (resp *types.ObjectResponse, err error) {
	object, err := l.svcCtx.ObjectsModel.FindOne(l.ctx, req.ObjId)
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

	return &types.ObjectResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.Object{
			ObjId:     object.ObjId,
			BucketId:  object.BucketId,
			UserId:    object.UserId,
			Key:       object.Key,
			Path:      object.Path,
			CreatedAt: object.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: object.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
