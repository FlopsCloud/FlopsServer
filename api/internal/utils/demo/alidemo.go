// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dyvmsapi "github.com/alibabacloud-go/dyvmsapi-20170525/v2/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

// 使用AK&SK初始化账号Client
func CreateDyvmsapiClient() (_result *dyvmsapi.Client, _err error) {
	accessKeyId := " "
	accessKeySecret := " "

	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	}
	_result = &dyvmsapi.Client{}
	_result, _err = dyvmsapi.NewClient(config)
	return _result, _err
}

func SendVerification(client *dyvmsapi.Client, verifyType *string, bizType *string, target *string) (_err error) {
	req := &dyvmsapi.SendVerificationRequest{
		VerifyType: verifyType,
		BizType:    bizType,
		Target:     target,
	}
	resp, _err := client.SendVerification(req)
	if _err != nil {
		return _err
	}

	if !tea.BoolValue(util.EqualString(resp.Body.Code, tea.String("OK"))) {
		_err = tea.NewSDKError(map[string]interface{}{
			"code":    tea.StringValue(resp.Body.Code),
			"message": tea.StringValue(resp.Body.Message),
		})
		return _err
	}

	console.Log(tea.String("------------sendVerification success-------------"))
	console.Log(util.ToJSONString(util.ToMap(resp)))
	return _err
}

func _main(args []*string) (_err error) {
	verifyType := args[0]
	bizType := args[1]
	target := args[2]
	client, _err := CreateDyvmsapiClient()
	if _err != nil {
		return _err
	}

	_err = SendVerification(client, verifyType, bizType, target)
	if _err != nil {
		return _err
	}
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
