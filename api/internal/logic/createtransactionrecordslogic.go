package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTransactionRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTransactionRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTransactionRecordsLogic {
	return &CreateTransactionRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTransactionRecordsLogic) CreateTransactionRecords(req *types.CreateTransactionRecordsRequest) response.Response {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		return response.Fail(response.UserNotExistCode, err.Error())
	}
	_, err = l.svcCtx.BalancesModel.FindOneByUserAndCurrency(l.ctx, req.UserId, req.CurrencyCode)
	if err != nil {
		return response.Fail(response.UserBalanceNotExistCode, err.Error())
	}
	_, err = l.svcCtx.TransactionRecordsModel.Insert(l.ctx, &model.TransactionRecords{
		Id:           0,
		UserId:       req.UserId,
		OrgId:        req.OrgId,
		TransType:    req.TransType,
		CurrencyCode: req.CurrencyCode,
		PayType:      req.PayType,
		Detail:       req.Detail,
		OrderNo:      req.OrderNo,
		Username:     user.Username,
		Amount:       req.Amount,
		Balance:      req.Balance,
	})
	if err != nil {
		response.Fail(response.ServerErrorCode, err.Error())
	}
	return response.OK("")
}
