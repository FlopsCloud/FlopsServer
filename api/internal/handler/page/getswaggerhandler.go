package page

import (
	"net/http"

	"fca/api/internal/logic/page"
	"fca/api/internal/svc"
	"fca/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Get swagger json
func GetSwaggerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := page.NewGetSwaggerLogic(r.Context(), svcCtx)
		resp, err := l.GetSwagger()
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call", err.Error()))
		} else {
			//w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(resp.Body))

		}
	}
}
