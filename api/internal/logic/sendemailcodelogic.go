package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/exp/rand"

	"gopkg.in/gomail.v2"
)

type SendemailcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendemailcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendemailcodeLogic {
	return &SendemailcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendemailcodeLogic) Sendemailcode(req *types.SendemailcodeRequest, clientIP string) (resp *types.Response, err error) {
	// Generate a random verification code
	code := generateVerificationCode()

	// Send the verification code to the provided email address
	err = l.sendVerificationEmail(req.Email, code)
	if err != nil {
		return nil, err
	}

	// Store the verification code and email in a database or cache
	err = l.storeVerificationData(req.Email, code, clientIP)
	if err != nil {
		return nil, err
	}

	resp = &types.Response{
		Code:    200,
		Message: "Verification code sent successfully",
	}

	return resp, nil
}

func generateVerificationCode() string {
	// Generate a random 6-digit verification code
	seed := uint64(time.Now().UnixNano())
	rand.Seed(seed)
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}

func (l *SendemailcodeLogic) sendVerificationEmail(email, code string) error {
	l.Logger.Infof("sendVerificationEmail start===========")

	m := gomail.NewMessage()
	m.SetHeader("From", "flopscloud@163.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email Verification Code")
	m.SetBody("text/html", fmt.Sprintf("<strong>Your verification code is: %s</strong>", code))

	// d := gomail.NewDialer("smtp.163.com", 587, "sender@example.com", "your-smtp-password")

	d := gomail.NewDialer("smtp.163.com", 25, "flopscloud@163.com", " ")

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}
	l.Logger.Infof("send mail success")

	return nil
}

func (l *SendemailcodeLogic) storeVerificationData(email, code, clientIP string) error {
	// Create a new VerificationCodes struct
	verificationCode := model.VerificationCodes{
		TargetType:     "email",
		TargetValue:    email,
		Code:           code,
		Ip:             clientIP,
		ExpirationTime: time.Now().Add(10 * time.Minute),
	}

	// Insert the verification code into the database
	_, err := l.svcCtx.VerificationCodesModel.Insert(l.ctx, &verificationCode)
	if err != nil {
		return fmt.Errorf("failed to store verification code: %v", err)
	}

	return nil
}
