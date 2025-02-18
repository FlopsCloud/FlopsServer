package logic

import (
	"context"
	"fca/common/response"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelloLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloLogic {
	return &HelloLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloLogic) Hello(req *types.HelloRequest) (resp *response.Response, err error) {
	// todo: add your logic here and delete this line
	return &response.Response{
		Code:    200,
		Message: "Hello " + req.Name,
	}, nil
}
