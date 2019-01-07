package email

import (
	"errors"

	"github.com/Timothylock/inventory-management/config"

	"gopkg.in/gomail.v2"
)

type Service struct {
	cfg    config.Config
	client *gomail.Dialer
}

func NewService(c config.Config) Service {
	if c.EmailSmtpServ == "" {
		return Service{
			cfg:    c,
			client: nil,
		}
	}

	return Service{
		cfg:    c,
		client: gomail.NewPlainDialer(c.EmailSmtpServ, c.EmailSmtpPort, c.EmailUsername, c.EmailPassword),
	}
}

func (s *Service) SendEmail(toEmail, subject, body string) error {
	if s.client == nil {
		return errors.New("your admin did not set up email properly. Please contact them to reset your password for you")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.EmailFromAddr)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.client.DialAndSend(m)
}
