package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServerTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServerTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServerTagLogic {
	return &ServerTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServerTagLogic) ServerTag() (resp *types.TagsResp, err error) {

	resp2 := new(types.TagsResp)
	resp2.Code = 0
	resp2.Message = "Success"

	// // 管理员权限验证 start
	// uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	// email, _ := l.ctx.Value("email").(string)

	// logx.Info("JWT uid=", uid, " Name=", email)
	// if uid != 1 {
	// 	resp.Code = 1
	// 	resp.Message = "Permission denied"
	// 	return resp, nil
	// }
	// // 管理员权限验证 end
	// l.Logger.Info(resp)
	tags, err := l.svcCtx.TagsModel.FindAll(l.ctx)
	// l.Logger.Info("==========")
	if err != nil {
		resp2.Code = 1
		resp2.Message = "Failed to fetch users: " + err.Error()
		return resp2, nil
	}

	l.Logger.Info(tags)

	tagList := make([]types.Tag, len(*tags))

	for i, tag := range *tags {

		tagList[i] = types.Tag{
			TagID:   int(tag.TagId),
			TagName: tag.TagName,
		}
	}

	resp2.Data = tagList

	return resp2, nil
}
