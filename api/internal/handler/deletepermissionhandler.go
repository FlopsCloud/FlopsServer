package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeletePermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeletePermissionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDeletePermissionLogic(r.Context(), svcCtx)
		resp := l.DeletePermission(req.Id)
		httpx.WriteJson(w, http.StatusOK, resp)
	}
}
