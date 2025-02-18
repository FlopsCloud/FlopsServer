package logic

import (
	"context"
	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/exp/rand"
	"gopkg.in/gomail.v2"
)

type OrgInviteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgInviteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgInviteUserLogic {
	return &OrgInviteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func generateInvaCode() string {
	// Generate a random 6-digit verification code
	seed := uint64(time.Now().UnixNano())
	rand.Seed(seed)
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}

func (l *OrgInviteUserLogic) sendInvitationEmail(email, code string) error {
	l.Logger.Infof("sendInvitationEmail start ")

	m := gomail.NewMessage()
	m.SetHeader("From", "flopscloud@163.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Invitation to join organization Code")
	m.SetBody("text/html", fmt.Sprintf("<strong>Your Invitation to join organization Code is : %s </strong>", code))

	// d := gomail.NewDialer("smtp.163.com", 587, "sender@example.com", "your-smtp-password")

	d := gomail.NewDialer("smtp.163.com", 25, "flopscloud@163.com", " ")

	if err := d.DialAndSend(m); err != nil {
		l.Logger.Errorf("failed to send Invitation email: %v", err)
		return fmt.Errorf("failed to send verification email: %v", err)
	}
	l.Logger.Infof("send Invitation mail success %s", email)

	return nil
}

func (l *OrgInviteUserLogic) OrgInviteUser(req *types.OrgInviteUserReq) (resp *types.Response, err error) {
	// Validate input
	if req.OrgId == 0 || req.Email == "" || req.Role == "" {
		return &types.Response{
			Code:    400,
			Message: "Invalid input: OrgId, Email, and Role are required",
		}, nil
	}

	// Validate role
	if req.Role != "member" && req.Role != "admin" {
		return &types.Response{
			Code:    400,
			Message: "Invalid role: must be either 'member' or 'admin'",
		}, nil
	}

	// Check if the organization exists
	org, err := l.svcCtx.OrganizationModel.FindOne(l.ctx, req.OrgId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.Response{
				Code:    404,
				Message: "Organization not found",
			}, nil
		}
		return nil, err
	}

	var user *model.Users
	if req.UserId != 0 {
		// If UserId is provided, verify that the email matches the user's email
		user, err = l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
		if err != nil {
			if err == model.ErrNotFound {
				return &types.Response{
					Code:    404,
					Message: "User not found",
				}, nil
			}
			return nil, err
		}
		if user.Email != req.Email {
			return &types.Response{
				Code:    400,
				Message: "Provided email does not match the user's email",
			}, nil
		}
	} else {
		// If UserId is not provided, find the user by email
		user, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
		if err != nil {
			if err == model.ErrNotFound {
				return &types.Response{
					Code:    404,
					Message: "User not found",
				}, nil
			}
			return nil, err
		}
	}

	// Since we don't have an OrgsUsersModel, we can't directly check if the user is already a member
	// or add them to the organization. Instead, we'll log this information and return a success message.

	// Log the invitation
	l.Logger.Infof("Invitation for user %s (ID: %d) to join organization %s (ID: %d) with role %s", user.Email, user.UserId, org.OrgName, org.OrgId, req.Role)

	// In a real implementation, you would typically:
	// 1. Check if the user is already a member of the organization
	// 2. Create an invitation record in the database
	// 3. Send an email to the user with the invitation

	code := generateInvaCode()

	_, err = l.svcCtx.InvitationModel.Insert(l.ctx, &model.Invitation{
		InviterId:       req.UserId,
		OrgId:           req.OrgId,
		Role:            req.Role,
		Status:          "pending",
		InviteeEmail:    req.Email,
		InvitationToken: code,
	})
	if err != nil {
		return &types.Response{
			Code:    500,
			Message: "Invitation fail",
			Info:    err.Error(),
		}, nil
	}

	err = l.sendInvitationEmail(req.Email, code)

	if err != nil {
		return &types.Response{
			Code:    500,
			Message: "Invitation sendInvitationEmail fail",
			Info:    err.Error(),
		}, nil
	}
	// _, err = l.svcCtx.OrgsUsersModel.Insert(l.ctx, &model.OrgsUsers{
	// 	UserId: req.UserId,
	// 	OrgId:  req.OrgId,
	// 	Role:   req.Role,
	// })

	// if err != nil {
	// 	return &types.Response{
	// 		Code:    500,
	// 		Message: "Invitation fail",
	// 		Info:    err.Error(),
	// 	}, nil
	// }

	return &types.Response{
		Code:    200,
		Message: "Invitation logged successfully. Implementation for adding user to organization is pending.",
	}, nil
}
