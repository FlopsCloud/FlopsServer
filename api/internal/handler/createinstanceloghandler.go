package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateInstanceLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateInstanceLogRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateInstanceLogLogic(r.Context(), svcCtx)
		resp, err := l.CreateInstanceLog(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
