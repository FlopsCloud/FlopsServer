package logic

import (
	"context"
	"encoding/json"
	"fca/api/internal/types"
	"fca/common/response"
	"fmt"

	"fca/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserBalanceLogic {
	return &GetUserBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserBalanceLogic) GetUserBalance() response.Response {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	condition := fmt.Sprintf("where user_id = %d", uid)
	items, err := l.svcCtx.BalancesModel.FindList(l.ctx, condition)
	if err != nil {
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	list := *items
	data := make([]types.UserBalanceData, len(list))
	for i := 0; i < len(list); i++ {
		data[i] = types.UserBalanceData{
			UserId:       list[i].UserId,
			Balance:      list[i].Balance,
			CurrencyCode: list[i].CurrencyCode,
		}
	}
	return response.OK(data)
}
