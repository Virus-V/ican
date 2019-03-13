package main

import (
	"crypto/tls"
	"ican/types"

	"go.uber.org/zap"

	gomail "gopkg.in/gomail.v2"
)

type sendmail struct {
	logger       *zap.Logger
	From         string
	FromName     string
	SMTPAddr     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

var _ types.SendMailService = (*sendmail)(nil)

func (s *sendmail) SendTo(to types.MailAddr, subject, body string) error {

	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", s.From, s.FromName)
	// 收件人
	m.SetAddressHeader("To", to.Addr, to.Name)
	// 主题
	m.SetHeader("Subject", subject)
	// 正文
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.SMTPAddr, s.SMTPPort, s.SMTPUsername, s.SMTPPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
