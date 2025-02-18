package ticket

import (
	"net/http"

	"fca/api/internal/logic/ticket"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// admin close ticket
func AdminCloseSupportTicketsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CloseSupportTicketRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ticket.NewAdminCloseSupportTicketsLogic(r.Context(), svcCtx)
		resp, err := l.AdminCloseSupportTickets(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
