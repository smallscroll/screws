package screws

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

//MailSender ...
type MailSender struct {
	Host     string //服务器
	Port     int    //端口
	Username string //用户名
	Password string //密码
	Form     string //发送者
	To       string //接收者
	Subject  string //主题
	Content  string //内容
}

//SendWithTLS ...
func (m *MailSender) SendWithTLS() error {

	//data
	header := make(map[string]string)
	header["From"] = m.Form + "<" + m.Username + ">"
	header["To"] = m.To
	header["Subject"] = m.Subject
	header["Content-Type"] = "text/html;chartset=UTF-8"
	body := m.Content

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	msg += "\r\n" + body

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	err := mailSendWithTLS(fmt.Sprintf("%s:%d", m.Host, m.Port), auth, m.Username, []string{m.To}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

//mailSendWithTLS ...
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
