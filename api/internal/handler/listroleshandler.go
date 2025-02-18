package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewListRolesLogic(r.Context(), svcCtx)
		resp := l.ListRoles(&req)
		httpx.WriteJson(w, http.StatusOK, resp)
	}
}
