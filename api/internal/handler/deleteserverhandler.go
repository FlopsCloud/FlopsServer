package handler

import (
	"net/http"

	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteServerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteServerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDeleteServerLogic(r.Context(), svcCtx)
		resp, err := l.DeleteServer(&req)
		if err != nil {
			//httpx.ErrorCtx(r.Context(), w, err)
			var ErrResp response.Response
			ErrResp.Code = response.ServerErrorCode
			ErrResp.Message = "Delete Server fail"
			ErrResp.Info = err.Error()
			httpx.OkJsonCtx(r.Context(), w, ErrResp)

		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
