package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListObjTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListObjTreeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewListObjTreeLogic(r.Context(), svcCtx)
		resp, err := l.ListObjTree(&req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    response.ServerErrorCode,
				Info:    err.Error(),
				Message: "ListObjTree Error",
			})
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
