package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CallbackRechargeAlipayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCallbackRechargeAlipayLogic(r.Context(), svcCtx)
		resp, err := l.CallbackRechargeAlipay()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
