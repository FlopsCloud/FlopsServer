package handler

import (
	"fmt"
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GoogleCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse request parameters from query or request body
		// example: id := r.URL.Query().Get("id")
		// example: var req types.CallbackRequest
		//          if err := httpx.Parse(r, &req); err != nil {
		//              httpx.ErrorCtx(r.Context(), w, err)
		//              return
		//          }

		code := r.URL.Query().Get("code")
		fmt.Println("GoogleCallbackHandler code: " + code)
		var req types.CallbackRequest
		req.Code = code

		l := logic.NewGoogleCallbackLogic(r.Context(), svcCtx)
		resp, err := l.GoogleCallback(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
