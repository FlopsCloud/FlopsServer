package logic

import (
	"context"
	"encoding/json"
	"fca/common/response"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServerListLogic {
	return &ServerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServerListLogic) ServerList(req *types.ServerListRequest) response.Response {
	//condition := fmt.Sprintf("where billing_method = %d and processor = %d", req.BillingMethod, req.Processor)

	// RegionId  int     `json:"regionId"`
	// TagId     int     `json:"tagId"`
	// IsPayMin 	int   `json:"isPayMin"`
	// IsPayDaily int  	`json:"isPayDaily"`
	// IsPayMonthly int `json:"isPayMonthly"`
	// IsPayYearly int  `json:"isPayYearly"`
	// GpuModel  int    `json:"gpuModel"`
	// GpuCount  int    `json:"gpuModel"`
	condition := fmt.Sprint("where 1=1")

	if req.RegionId > 0 {
		condition += fmt.Sprintf(" and region_id = %d", req.RegionId)
	}
	if req.TagId > 0 {
		condition += fmt.Sprintf(" and server_id in (select server_id from server_tags where tag_id= %d )", req.TagId)
	}
	if req.IsPayMin > 0 {
		condition += fmt.Sprintf(" and is_pay_min = %d", req.IsPayMin)
	}
	if req.IsPayDaily > 0 {
		condition += fmt.Sprintf(" and is_pay_daily = %d", req.IsPayDaily)
	}
	if req.IsPayMonthly > 0 {
		condition += fmt.Sprintf(" and is_pay_monthly = %d", req.IsPayMonthly)
	}
	if req.IsPayYearly > 0 {
		condition += fmt.Sprintf(" and is_pay_yearly = %d", req.IsPayYearly)
	}

	if req.GpuCount > 0 {
		condition += fmt.Sprintf(" and gpu_count = %d", req.GpuCount)
	}
	if req.GpuModel > 0 {
		condition += fmt.Sprintf(" and gpu_model = %d", req.GpuModel)
	}
	// if req.GraphicsMemory > 0 {
	// 	condition += fmt.Sprintf(" and gpu_mem = %d", req.GraphicsMemory)
	// }
	// if req.CpuModel > 0 {
	// 	condition += fmt.Sprintf(" and cpu_model = %d", req.CpuModel)
	// }
	// if req.GpuCount > 0 {
	// 	condition += fmt.Sprintf(" and thread_count = %d", req.GpuCount)
	// }

	total, err := l.svcCtx.ServerModel.Count(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	condition += fmt.Sprintf(" limit %d,%d", (req.Page-1)*req.PageSize, req.PageSize)

	data, err := l.svcCtx.ServerModel.FindList(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	var serviceList []types.Server
	err = json.Unmarshal(bytes, &serviceList)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	return response.OK(&types.ServerListDataResp{Servers: serviceList, Total: total})
}
