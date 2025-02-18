package web

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWebslideLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWebslideLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWebslideLogic {
	return &UpdateWebslideLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWebslideLogic) UpdateWebslide(req *types.WebslideRespData) (resp *types.Response, err error) {
	resp = new(types.Response)

	result, err := l.svcCtx.WebConfigModel.FindOne(l.ctx, 1)
	var webConfigObj WebConfigObj
	if err != nil {

		resp.Code = response.ServerErrorCode
		resp.Info = err.Error()
		resp.Message = "Error while get webslide"
		return resp, nil
	}
	err = json.Unmarshal([]byte(result.Config), &webConfigObj)
	if err != nil {

		resp.Code = response.ServerErrorCode
		resp.Info = err.Error()
		resp.Message = "Error while get webslide"
		return resp, nil
	}
	webConfigObj.Images = req.Images
	data, err := json.Marshal(webConfigObj)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "Error while update webslide"
		resp.Info = err.Error()
		return resp, nil
	}

	l.svcCtx.WebConfigModel.Update(l.ctx, &model.WebConfig{
		Id:     1,
		Config: string(data),
	})

	resp.Code = response.SuccessCode
	resp.Message = "OK"

	return resp, nil

}
