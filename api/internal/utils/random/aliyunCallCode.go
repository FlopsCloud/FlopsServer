// This file is auto-generated, don't edit it. Thanks.
package random

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dyvmsapi "github.com/alibabacloud-go/dyvmsapi-20170525/v2/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

// DyvmsapiClientInterface defines the methods we need from the Dyvmsapi client
type DyvmsapiClientInterface interface {
	SendVerification(request *dyvmsapi.SendVerificationRequest) (*dyvmsapi.SendVerificationResponse, error)
}

// CreateDyvmsapiClient initializes and returns a new Dyvmsapi client
func CreateDyvmsapiClient() (DyvmsapiClientInterface, error) {
	accessKeyId := " "
	accessKeySecret := " "

	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	}
	return dyvmsapi.NewClient(config)
}

// SendVerification sends a verification request
func SendVerification(client DyvmsapiClientInterface, verifyType, bizType, target *string) error {
	req := &dyvmsapi.SendVerificationRequest{
		VerifyType: verifyType,
		BizType:    bizType,
		Target:     target,
	}

	fmt.Print(req)
	resp, err := client.SendVerification(req)
	if err != nil {
		return err
	}

	if !tea.BoolValue(util.EqualString(resp.Body.Code, tea.String("OK"))) {
		return tea.NewSDKError(map[string]interface{}{
			"code":    tea.StringValue(resp.Body.Code),
			"message": tea.StringValue(resp.Body.Message),
		})
	}

	console.Log(tea.String("------------sendVerification success-------------"))
	console.Log(util.ToJSONString(util.ToMap(resp)))
	return nil
}

// SendVerificationCode sends a verification code to the specified phone number
func SendVerificationCode(phone string, client DyvmsapiClientInterface) error {
	if client == nil {
		var err error
		client, err = CreateDyvmsapiClient()
		if err != nil {
			return err
		}
	}

	verifyType := "SMS"
	bizType := "CONTACT"
	target := phone
	return SendVerification(client, &verifyType, &bizType, &target)
}
