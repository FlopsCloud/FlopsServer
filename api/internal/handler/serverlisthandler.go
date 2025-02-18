package handler

import (
	"fca/common/response"
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ServerListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServerListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, response.Fail(response.InvalidRequestParamCode, err.Error()))
			return
		}

		l := logic.NewServerListLogic(r.Context(), svcCtx)
		httpx.OkJsonCtx(r.Context(), w, l.ServerList(&req))
	}
}
