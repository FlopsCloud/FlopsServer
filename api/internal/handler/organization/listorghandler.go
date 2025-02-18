package organization

import (
	"net/http"

	"fca/api/internal/logic/organization"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListOrgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListOrgRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call", err.Error()))
			return
		}

		l := organization.NewListOrgLogic(r.Context(), svcCtx)
		resp, err := l.ListOrg(&req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, response.FailWithInfo(response.InvalidRequestParamCodeInHandler, "Error while handle call", err.Error()))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
