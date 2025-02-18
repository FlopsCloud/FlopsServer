package logic

import (
	"context"
	"encoding/json"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateResourceLogic {
	return &CreateResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateResourceLogic) CreateResource(req *types.ResourceCreateRequest) (resp *types.ResourceResponse, err error) {
	// Create new resource

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	email, _ := l.ctx.Value("email").(string)

	logx.Info("JWT uid=", uid, " Name=", email)

	resource := &model.Resources{
		ResourceType:  req.ResourceType,
		UnitMinPrice:  req.UnitMinPrice,
		UnitHourPrice: req.UnitHourPrice,
		UnitDayPrice:  req.UnitDayPrice,
		IsDeleted:     0,
		CreatedAt:     time.Now(),
		CreatedBy:     uint64(uid),
	}

	result, err := l.svcCtx.ResourcesModel.Insert(l.ctx, resource)
	if err != nil {
		return nil, err
	}

	resourceId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get the created resource
	createdResource, err := l.svcCtx.ResourcesModel.FindOne(l.ctx, uint64(resourceId))
	if err != nil {
		return nil, err
	}

	return &types.ResourceResponse{
		Response: types.Response{
			Code:    0,
			Message: "success",
		},
		Data: types.ResourceItem{
			ResourceId: createdResource.ResourceId,

			ResourceType:  createdResource.ResourceType,
			UnitMinPrice:  createdResource.UnitMinPrice,
			UnitHourPrice: createdResource.UnitHourPrice,
			UnitDayPrice:  createdResource.UnitDayPrice,
			IsDeleted:     createdResource.IsDeleted,
			CreatedAt:     createdResource.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     createdResource.CreatedBy,
		},
	}, nil
}
