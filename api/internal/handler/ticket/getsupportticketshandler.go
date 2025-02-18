package ticket

import (
	"net/http"

	"fca/api/internal/logic/ticket"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSupportTicketsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSupportTicketsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ticket.NewGetSupportTicketsLogic(r.Context(), svcCtx)
		resp, err := l.GetSupportTickets(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
