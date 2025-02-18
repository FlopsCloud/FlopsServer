package logic

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/volcengine/volc-sdk-golang/service/sms"
	"github.com/zeromicro/go-zero/core/logx"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

type SendphonecodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendphonecodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendphonecodeLogic {
	return &SendphonecodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendphonecodeLogic) Sendphonecode(req *types.SendphonecodeRequest, ip string) (resp *types.Response, err error) {
	// Generate a random 6-digit code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	l.Logger.Infof("Generated verification code: %s for phone number: %s", code, req.Phonenum)

	// Save the code to the database
	verificationCode := &model.VerificationCodes{
		TargetType:     "phone",
		TargetValue:    req.Phonenum,
		Code:           code,
		Ip:             ip,
		ExpirationTime: time.Now().Add(5 * time.Minute),
	}
	_, err = l.svcCtx.VerificationCodesModel.Insert(l.ctx, verificationCode)
	if err != nil {
		l.Logger.Errorf("Failed to save verification code: %v", err)
		return &types.Response{
			Code:    500,
			Message: "Failed to save verification code",
		}, err
	}
	l.Logger.Info("Verification code saved to database")

	// Choose SMS service based on configuration
	useAlibabaCloud := true // Set this based on your configuration

	var smsErr error
	if useAlibabaCloud {
		smsErr = l.sendSMSAlibaba(req.Phonenum, code)
	} else {
		smsErr = l.sendSMSVolc(req.Phonenum, code)
	}

	if smsErr != nil {
		l.Logger.Errorf("Failed to send SMS: %v", smsErr)
		return &types.Response{
			Code:    500,
			Message: "Failed to send SMS",
			Info:    fmt.Sprintf("SMS send error: %s", smsErr.Error()),
		}, smsErr
	}

	l.Logger.Infof("SMS sent successfully to %s", req.Phonenum)
	return &types.Response{
		Code:    200,
		Message: "Verification code sent successfully",
	}, nil
}

func (l *SendphonecodeLogic) sendSMSVolc(phoneNumber, code string) error {
	testAk, testSk := " ", " =="
	sms.DefaultInstance.Client.SetAccessKey(testAk)
	sms.DefaultInstance.Client.SetSecretKey(testSk)
	reqsms := &sms.SmsRequest{
		SmsAccount:    "7e2a96f5", //消息组ID
		Sign:          "短信服务",
		TemplateID:    "SPT_09a29a26", //ST_7e335c5b
		TemplateParam: fmt.Sprintf(`{"code":"%s"}`, code),
		PhoneNumbers:  phoneNumber,
		Tag:           "tag",
	}
	smsResponse, _, err := sms.DefaultInstance.Send(reqsms)
	if err != nil {
		return err
	}

	if smsResponse.ResponseMetadata.Error != nil {
		return fmt.Errorf("SMS service returned an error: %s", smsResponse.ResponseMetadata.Error.Message)
	}

	return nil
}

func (l *SendphonecodeLogic) sendSMSAlibaba(phoneNumber, code string) error {
	// Replace these with your actual Alibaba Cloud access key and secret
	accessKeyId := " "
	accessKeySecret := " "

	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	}

	client, err := dysmsapi.NewClient(config)
	if err != nil {
		return err
	}

	sendSmsRequest := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumber),
		SignName:      tea.String("广州丰捷软件技术"),
		TemplateCode:  tea.String("SMS_297980012"),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
	}

	sendresponse, err := client.SendSms(sendSmsRequest)
	if err != nil {
		return err
	}

	if *sendresponse.Body.Code != "OK" {
		return fmt.Errorf("SMS send failed with code: %s, message: %s", *sendresponse.Body.Code, *sendresponse.Body.Message)
	}

	//
	// 	"statusCode": 200,
	//    "body": {
	//       "Code": "isv.AMOUNT_NOT_ENOUGH",
	//       "Message": "账户余额不足",
	//       "RequestId": "2E3C074F-3C9F-5664-B09F-3086CC8769B7"
	//    }

	// "statusCode": 200,
	// "body": {
	//    "BizId": "645924827285632285^0",
	//    "Code": "OK",
	//    "Message": "OK",
	//    "RequestId": "6BB40AEE-D820-5F61-88FA-2607D7C1EB6A"
	// }

	l.Logger.Infof("resp", sendresponse)

	return nil
}
