package admin

import (
	"net/http"

	"fca/api/internal/logic/admin"
	"fca/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminPanelInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminPanelInfoLogic(r.Context(), svcCtx)
		resp, err := l.AdminPanelInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
