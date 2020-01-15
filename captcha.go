package screws

import "fmt"

//ICaptcha 验证码管理器接口
type ICaptcha interface {
	Send(accountType, to, from, subject string, expiration int32) error
	Get(account string) (string, error)
	Expiration(account string) (int32, error)
}

//captcha 验证码管理器
type captcha struct {
	MailSender IMailSender
	SmsSender  IAlisms
	Cache      ICache
}

//NewCaptcha 初始化验证码管理器
func NewCaptcha(mailSender IMailSender, smsSender IAlisms, cache ICache) ICaptcha {
	return &captcha{
		MailSender: mailSender,
		SmsSender:  smsSender,
		Cache:      cache,
	}
}

//Send 验证码发送: mail/mobile
func (c *captcha) Send(accountType, to string, from, subject string, expiration int32) error {
	code := NewTinyTools().DigitalCaptcha()
	switch accountType {
	case "mail":
		if from == "" {
			from = "System"
		}
		if subject == "" {
			subject = "Captcha"
		}
		if err := c.MailSender.SendWithTLS(from, to, subject, fmt.Sprintf("您的验证码是： %s，%d分钟内有效。", code, expiration/60)); err != nil {
			return err
		}
	case "mobile":
		if err := c.SmsSender.SendCaptcha(to, code); err != nil {
			return err
		}
	}
	if err := c.Cache.Set(to, []byte(code), 0, expiration); err != nil {
		return err
	}
	return nil
}

//Get 验证码获取
func (c *captcha) Get(account string) (string, error) {
	v, err := c.Cache.Get(account)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

//Get 有效期获取
func (c *captcha) Expiration(account string) (int32, error) {
	v, err := c.Cache.Expiration(account)
	if err != nil {
		return v, err
	}
	return v, nil
}

//Delete 验证码删除
func (c *captcha) Delete(account string) error {
	return c.Cache.Delete(account)
}