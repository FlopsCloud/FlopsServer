package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListK8sServerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListK8sServerLogic(r.Context(), svcCtx)
		jwtToken := r.Header.Get("Authorization")
		resp, err := l.ListK8sServer(jwtToken)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
