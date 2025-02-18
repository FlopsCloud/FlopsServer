package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req types.DeleteResourceRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// resourceId := r.PathValue("resourceId")
		// id, err := strconv.ParseUint(resourceId, 10, 64)
		// if err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// 	return
		// }

		l := logic.NewDeleteResourceLogic(r.Context(), svcCtx)
		resp, err := l.DeleteResource(req.Id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
