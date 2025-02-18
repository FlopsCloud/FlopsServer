package logic

import (
	"context"
	"fca/common/response"
	"fca/model"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBaseDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBaseDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseDataLogic {
	return &UpdateBaseDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBaseDataLogic) UpdateBaseData(req *types.BaseData) response.Response {
	// 查一下看有没有重复的type和name

	existId, err := l.svcCtx.BaseDataModel.FindByTypeName(l.ctx, req.DataType, req.Name)

	if existId > 0 {
		return response.Fail(response.InvalidRequestParamCode, fmt.Sprintf("类型为%s的数据已包含值: %s，不能重复", req.DataType, req.Name))
	}

	// 没有，那么就根据req.id 来更新这条记录

	err = l.svcCtx.BaseDataModel.Update(l.ctx, &model.BaseData{
		Id:       req.Id,
		Name:     req.Name,
		DataType: req.DataType,
		DataId:   req.DataID,
		Comment:  req.Comment,
	})
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	l.svcCtx.RedisClient.Del("BaseDataFilter")

	return response.OK(&types.BaseData{
		Id:       req.Id,
		Name:     req.Name,
		DataID:   req.DataID,
		DataType: req.DataType,
		Comment:  req.Comment,
	})
}
