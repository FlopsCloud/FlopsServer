package page

import (
	"net/http"

	"fca/api/internal/logic/page"
	"fca/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Get Privacy policy
func GetPPHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := page.NewGetPPLogic(r.Context(), svcCtx)
		resp, err := l.GetPP()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(resp.Body))

		}
	}
}
