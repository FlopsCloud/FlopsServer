package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListObjectLogic {
	return &ListObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListObjectLogic) ListObject(req *types.ListObjectRequest) (resp *types.ListObjectResponse, err error) {
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
	total, err := l.svcCtx.ObjectsModel.Count(l.ctx, req.BucketId, req.Key)
	if err != nil {
		return &types.ListObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get total count",
				Info:    err.Error(),
			},
		}, nil
	}

	// Get objects with pagination
	objects, err := l.svcCtx.ObjectsModel.FindMany(l.ctx, req.BucketId, req.Key, offset, req.PageSize)
	if err != nil {
		return &types.ListObjectResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get objects",
				Info:    err.Error(),
			},
		}, nil
	}

	// Convert to response type
	var objectList []types.Object
	for _, obj := range objects {
		objectList = append(objectList, types.Object{
			ObjId:     obj.ObjId,
			BucketId:  obj.BucketId,
			UserId:    obj.UserId,
			Key:       obj.Key,
			Path:      obj.Path,
			CreatedAt: obj.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: obj.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ListObjectResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ListObjectResponseData{
			Objects: objectList,
			Total:   total,
		},
	}, nil
}
