package screws

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

//IAlisms 阿里短信接口
type IAlisms interface {
	Send(phoneNumbers, templateParam string) error
}

//NewAlisms 初始化阿里短信
func NewAlisms(accessKeyID, accessKeySecret, signName, templateCode string) IAlisms {
	return &alismsSender{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		SignName:        signName,
		TemplateCode:    templateCode,
	}
}

//alismsSender 阿里短信
type alismsSender struct {
	AccessKeyID     string //AccessKeyID
	AccessKeySecret string //AccessKeySecret
	SignName        string //短信签名
	TemplateCode    string //模板编号
}

//alismsReply 调用返回
type alismsReply struct {
	Message   string
	RequestID string
	BizID     string
	Code      string
}

//Send 发送
func (as *alismsSender) Send(phoneNumbers, templateParam string) error {
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
	request.QueryParams["PhoneNumbers"] = phoneNumbers
	request.QueryParams["SignName"] = as.SignName
	request.QueryParams["TemplateCode"] = as.TemplateCode
	request.QueryParams["TemplateParam"] = templateParam

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
