package web

import (
	"context"
	"encoding/json"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWebslideLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWebslideLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWebslideLogic {
	return &ListWebslideLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type WebConfigObj struct {
	Images []string `json:"images"`
}

func (l *ListWebslideLogic) ListWebslide() (resp *types.WebslideResp, err error) {

	resp = new(types.WebslideResp)
	resp.Code = 200
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

	resp.Code = response.SuccessCode
	resp.Message = "OK"
	resp.Data = new(types.WebslideRespData)
	resp.Data.Images = webConfigObj.Images
	return resp, nil

}
