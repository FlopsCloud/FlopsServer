package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBaseDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBaseDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBaseDataLogic {
	return &DeleteBaseDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBaseDataLogic) DeleteBaseData(req *types.DeleteBaseDataRequest) (resp *types.Response, err error) {

	err = l.svcCtx.BaseDataModel.Delete(l.ctx, req.Id)
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: "删除失败",
		}, nil
	}
	l.svcCtx.RedisClient.Del("BaseDataFilter")
	return &types.Response{
		Code:    200,
		Message: "删除成功",
	}, nil

}
