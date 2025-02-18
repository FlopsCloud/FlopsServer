package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListPermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PermissionListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewListPermissionsLogic(r.Context(), svcCtx)
		resp := l.ListPermissions(&req)
		httpx.WriteJson(w, http.StatusOK, resp)
	}
}
