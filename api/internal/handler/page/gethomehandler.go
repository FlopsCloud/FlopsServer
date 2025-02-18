package page

import (
	"net/http"

	"fca/api/internal/logic/page"
	"fca/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Get home page
func GetHomeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := page.NewGetHomeLogic(r.Context(), svcCtx)
		resp, err := l.GetHome()

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(resp.Body))

			// httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
