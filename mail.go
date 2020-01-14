package screws

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

//IMailSender 邮件发送器接口
type IMailSender interface {
	SendWithTLS(from, to, subject, content string) error
}

//mailSender 邮件发送器
type mailSender struct {
	Host     string //服务器
	Port     string //端口
	Username string //用户名
	Password string //密码
}

//NewMailSender 初始化邮件发送器(邮件服务器，端口，用户名，密码)
func NewMailSender(host, port, username, password string) IMailSender {
	return &mailSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

//SendWithTLS 加密发送(发送者，接收者，主题，内容)
func (m *mailSender) SendWithTLS(from, to, subject, content string) error {

	header := make(map[string]string)
	header["From"] = from + "<" + m.Username + ">"
	header["To"] = to
	header["Subject"] = subject
	header["Content-Type"] = "text/html;chartset=UTF-8"
	body := content

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	msg += "\r\n" + body

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	err := mailSendWithTLS(fmt.Sprintf("%s:%s", m.Host, m.Port), auth, m.Username, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func mailSendWithTLS(addr string, a smtp.Auth, from string, to []string, msg []byte) error {

	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	host, _, _ := net.SplitHostPort(addr)
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	err = c.Quit()
	if err != nil {
		return err
	}
	return nil
}
