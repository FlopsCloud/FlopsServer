package wx

import (
	"net/http"

	"fca/api/internal/logic/wx"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GenerateQRCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateQRCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call", err.Error()))
			return
		}

		l := wx.NewGenerateQRCodeLogic(r.Context(), svcCtx)
		resp, err := l.GenerateQRCode(&req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call", err.Error()))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
