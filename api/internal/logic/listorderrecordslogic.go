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

type ListOrderRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrderRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrderRecordsLogic {
	return &ListOrderRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOrderRecordsLogic) ListOrderRecords(req *types.OrderRecordsListRequest) response.Response {
	condition := "where 1=1"
	if req.UserId > 0 {
		condition += fmt.Sprintf(" and user_id = %d", req.UserId)
	}
	if req.OrderType > 0 {
		condition += fmt.Sprintf(" and order_type = %d", req.OrderType)
	}
	if req.OrderNo != "" {
		condition += fmt.Sprintf(" and order_no = '%s'", req.OrderNo)
	}
	if req.BeginCreatedAt > 0 {
		condition += fmt.Sprintf(" and created_at > %d", req.BeginCreatedAt)
	}
	if req.EndCreatedAt > 0 {
		condition += fmt.Sprintf(" and created_at < %d", req.EndCreatedAt)
	}

	total, err := l.svcCtx.OrderRecordsModel.Count(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	condition += fmt.Sprintf(" limit %d,%d", (req.Page-1)*req.PageSize, req.PageSize)

	data, err := l.svcCtx.OrderRecordsModel.FindList(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	var items []types.TransactionRecords
	err = json.Unmarshal(bytes, &items)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	return response.OK(&types.TransactionRecordsListResponseData{Items: items, Total: total})
}
