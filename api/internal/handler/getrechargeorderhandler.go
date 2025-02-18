package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRechargeOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRechargeOrderRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call", err.Error()))
			return
		}

		l := logic.NewGetRechargeOrderLogic(r.Context(), svcCtx)
		httpx.OkJsonCtx(r.Context(), w, l.GetRechargeOrder(&req))
	}
}
