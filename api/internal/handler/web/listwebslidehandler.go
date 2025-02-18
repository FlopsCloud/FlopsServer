package web

import (
	"net/http"

	"fca/api/internal/logic/web"
	"fca/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListWebslideHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := web.NewListWebslideLogic(r.Context(), svcCtx)
		resp, err := l.ListWebslide()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
