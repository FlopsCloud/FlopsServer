package handler

import (
	"errors"
	"fca/common/response"
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CallbackRechargeWeixinHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCallbackRechargeWeixinLogic(r.Context(), svcCtx, r)
		resp := l.CallbackRechargeWeixin()
		if resp.Code != response.SuccessCode {
			httpx.ErrorCtx(r.Context(), w, errors.New(resp.Message))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
