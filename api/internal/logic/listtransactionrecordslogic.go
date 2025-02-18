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

type ListTransactionRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTransactionRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTransactionRecordsLogic {
	return &ListTransactionRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTransactionRecordsLogic) ListTransactionRecords(req *types.TransactionRecordsListRequest) response.Response {
	condition := "where 1=1"
	if req.UserId > 0 {
		condition += fmt.Sprintf(" and user_id = %d", req.UserId)
	}
	if req.OrgId > 0 {
		condition += fmt.Sprintf(" and org_id = %d", req.OrgId)
	}
	if req.TransType > 0 {
		condition += fmt.Sprintf(" and trans_type = %d", req.TransType)
	}
	if req.CurrencyCode != "" {
		condition += fmt.Sprintf(" and currency_code = '%s'", req.CurrencyCode)
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

	total, err := l.svcCtx.TransactionRecordsModel.Count(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	condition += fmt.Sprintf(" limit %d,%d", (req.Page-1)*req.PageSize, req.PageSize)

	data, err := l.svcCtx.TransactionRecordsModel.FindList(l.ctx, condition)
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
