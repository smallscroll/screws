package screws

import (
	"errors"
	"fmt"
)

//ICaptcha 验证码管理器接口
type ICaptcha interface {
	Send(accountType, to, from, subject string, expiration int32) error
	Get(account string) (string, error)
	Expiration(account string) (int32, error)
	Delete(account string) error
}

//captcha 验证码管理器
type captcha struct {
	MailSender IMailSender
	SmsSender  IAlisms
	Cache      ICache
}

//NewCaptcha 初始化验证码管理器
func NewCaptcha(captchaMailSender IMailSender, captchaSmsSender IAlisms, cache ICache) (ICaptcha, error) {
	if cache == nil {
		return nil, errors.New("Cache instance is nil")
	}
	return &captcha{
		MailSender: captchaMailSender,
		SmsSender:  captchaSmsSender,
		Cache:      cache,
	}, nil

}

//Send 验证码发送: mail/mobile
func (c *captcha) Send(accountType, to string, from, subject string, expiration int32) error {
	if c.Cache == nil {
		return errors.New("Cache instance is nil")
	}
	code := NewTinyTools().DigitalCaptcha()
	switch accountType {
	case "mail":
		if from == "" {
			from = "System"
		}
		if subject == "" {
			subject = "Captcha"
		}
		if expiration == 0 {
			expiration = 1800
		}
		if err := c.MailSender.SendWithTLS(from, to, subject, fmt.Sprintf("验证码：%s", code)); err != nil {
			return err
		}
	case "mobile":
		if err := c.SmsSender.Send(to, fmt.Sprintf(`{"code":"%s"}`, code)); err != nil {
			return err
		}
	default:
		return errors.New("Captcha send: accountType error")
	}
	if err := c.Cache.Set(to, []byte(code), 0, expiration); err != nil {
		return err
	}
	return nil
}

//Get 验证码获取
func (c *captcha) Get(account string) (string, error) {
	if c.Cache == nil {
		return "", errors.New("Cache instance is nil")
	}
	v, err := c.Cache.Get(account)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

//Get 有效期获取
func (c *captcha) Expiration(account string) (int32, error) {
	if c.Cache == nil {
		return -1, errors.New("Cache instance is nil")
	}
	v, err := c.Cache.Expiration(account)
	if err != nil {
		return v, err
	}
	return v, nil
}

//Delete 验证码删除
func (c *captcha) Delete(account string) error {
	if c.Cache == nil {
		return errors.New("Cache instance is nil")
	}
	return c.Cache.Delete(account)
}
