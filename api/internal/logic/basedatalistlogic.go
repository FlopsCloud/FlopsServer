package logic

import (
	"context"
	"fca/common/response"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseDataListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBaseDataListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseDataListLogic {
	return &BaseDataListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseDataListLogic) BaseDataList(req *types.BaseDataListRequest) response.Response {
	condition := ""
	if len(req.Name) > 0 && req.DataType != "" {
		condition += fmt.Sprintf(" where name like %s and data_type like  %s", "'%"+req.Name+"%'", "'%"+req.DataType+"%'")
	} else if len(req.Name) > 0 {
		condition += fmt.Sprintf(" where name like %s", "'%"+req.Name+"%'")
	} else if req.DataType != "" {
		condition += fmt.Sprintf(" where data_type like %s", "'%"+req.DataType+"%'")
	}

	condition += fmt.Sprintf(" limit %d,%d", (req.Page-1)*req.PageSize, req.PageSize)
	data, err := l.svcCtx.BaseDataModel.FindList(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	total, err := l.svcCtx.BaseDataModel.Count(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	//baseDataList := make([]types.BaseData, len(*data))
	baseDataList := make([]types.BaseData, 0)
	for _, v := range *data {
		baseDataList = append(baseDataList, types.BaseData{
			Id:       v.Id,
			DataID:   v.DataId,
			DataType: v.DataType,
			Name:     v.Name,
			Comment:  v.Comment,
		})

	}

	return response.OK(&types.BaseDataListResponseData{
		BaseDataList: baseDataList,
		Total:        uint64(total),
	})
}
