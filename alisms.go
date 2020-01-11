package screws

import (
	"encoding/json"
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

//AlismsSender ...
type AlismsSender struct {
	AccessKeyID     string //AccessKeyID
	AccessKeySecret string //AccessKeySecret
	SignName        string //短信签名
	TemplateCode    string //短信模板
	Receiver        string //短信接收者
	Content         string //短信内容
}

//alismsReply 接口返回
type alismsReply struct {
	Message   string
	RequestID string
	BizID     string
	Code      string
}

//SendCaptcha 验证码
func (as *AlismsSender) SendCaptcha() error {

	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", as.AccessKeyID, as.AccessKeySecret)
	if err != nil {
		return err
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = as.Receiver
	request.QueryParams["SignName"] = as.SignName
	request.QueryParams["TemplateCode"] = as.TemplateCode
	request.QueryParams["TemplateParam"] = "{'code':" + as.Content + "}"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}
	var reply alismsReply
	json.Unmarshal(response.GetHttpContentBytes(), &reply)
	if reply.Message != "OK" {
		return errors.New(reply.Message)
	}
	return nil
}
