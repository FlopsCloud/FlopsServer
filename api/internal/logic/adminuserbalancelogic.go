package logic

import (
	"context"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/api/internal/utils"
	"fca/common/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUserBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserBalanceLogic {
	return &AdminUserBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserBalanceLogic) AdminUserBalance(req *types.AdminUserBalanceRequest) (resp *types.AdminUserBalanceResponse, err error) {

	resp = &types.AdminUserBalanceResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
	}

	condition := fmt.Sprintf("where balances.user_id = %d and balances.currency_code = 'CNY'", req.UserId)
	balance, err := l.svcCtx.BalancesModel.FindListView(l.ctx, condition)
	if err != nil {
		resp.Code = response.ServerErrorCode
		resp.Message = "Failed to fetch user balance"
		return resp, err
	}

	userBalance := make([]types.UserBalance, len(*balance))
	for i, v := range *balance {

		var des types.UserBalance
		utils.CopyStruct(&v, &des, false)
		userBalance[i] = des
	}

	resp.Data = &types.AdminUserBalanceResponseData{
		UserBalance: userBalance,
	}

	return
}
