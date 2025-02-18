package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// User applies to join an organization
func ApplyToJoinOrglistHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyToJoinListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewApplyToJoinOrglistLogic(r.Context(), svcCtx)
		resp, err := l.ApplyToJoinOrglist(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
