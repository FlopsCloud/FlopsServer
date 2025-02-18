package logic

import (
	"context"
	"encoding/json"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseDataFilterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBaseDataFilterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseDataFilterLogic {
	return &BaseDataFilterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseDataFilterLogic) BaseDataFilter() (res response.Response) {

	cacheData, err := l.svcCtx.RedisClient.Get("BaseDataFilter")
	if err == nil {

		var data types.BaseDataFilterResponseData
		err = json.Unmarshal([]byte(cacheData), &data)
		if err == nil {
			return response.OK(data)
		}

	}

	data := &types.BaseDataFilterResponseData{
		BillingMethod:  transform(l, "计费方式"),
		Region:         transform(l, "地区"),
		Supplier:       transform(l, "供应商"),
		Processor:      transform(l, "处理器类型"),
		GPUModel:       transform(l, "GPU"),
		GraphicsMemory: transform(l, "显存"),
		CPUModel:       transform(l, "CPU"),
		ProcessorCount: transform(l, "核数"),
		Image:          transform(l, "镜像"),
		IP:             transform(l, "IP"),
	}

	// 转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return response.Error("转换为 JSON 失败")
	}

	// 保存到 Redis
	// l.svcCtx.RedisClient.Setex("BaseDataFilter", string(jsonData), 60)
	l.svcCtx.RedisClient.Set("BaseDataFilter", string(jsonData))

	return response.OK(data)

}

func transform(l *BaseDataFilterLogic, datatype string) (baseDataList []types.BaseData) {
	data, _ := l.svcCtx.BaseDataModel.FindByType(l.ctx, datatype)

	// baseDataList := make([]types.BaseData, len(*data))
	baseDataList = make([]types.BaseData, 0)
	l.Logger.Infof("baseDataList len %d", len(*data))
	for _, v := range *data {
		baseDataList = append(baseDataList, types.BaseData{
			Id:       v.Id,
			DataID:   v.DataId,
			DataType: v.DataType,
			Name:     v.Name,
			Comment:  v.Comment,
		})

	}

	return baseDataList
}
