package ticket

import (
	"context"
	"encoding/json"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplySupportTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplySupportTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplySupportTicketLogic {
	return &ReplySupportTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplySupportTicketLogic) ReplySupportTicket(req *types.ReplySupportTicketRequest) (resp *types.Response, err error) {

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	role, _ := l.ctx.Value("role").(string)

	ticket, err := l.svcCtx.SupportTicketsModel.FindOne(l.ctx, req.TicketId)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "找不到工单",
			Info:    err.Error(),
		}, nil
	}
	if ticket.UserId != uint64(uid) && (role != "superadmin" && role != "admin") {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "不能回复他人的工单",
		}, nil
	}

	if role == "superadmin" || role == "admin" {
		ticket.Status = "in-progress"
		err = l.svcCtx.SupportTicketsModel.Update(l.ctx, ticket)
		if err != nil {
			return &types.Response{
				Code:    response.ServerErrorCode,
				Message: "更新工单状态失败",
				Info:    err.Error(),
			}, nil
		}
	}

	rt := model.TicketReplies{
		TicketId: req.TicketId,
		Content:  req.Content,
		UserId:   uint64(uid),
		Images:   req.Images,
	}

	result, err := l.svcCtx.TicketRepliesModel.Insert(l.ctx, &rt)
	if err != nil {
		return &types.Response{
			Code:    response.ServerErrorCode,
			Message: "回复错误失败",
			Info:    err.Error(),
		}, nil
	}
	replyId, _ := result.LastInsertId()

	return &types.Response{
		Code:    response.SuccessCode,
		Message: "回复成功",
		Info:    fmt.Sprintf("%d", replyId),
	}, nil
}
