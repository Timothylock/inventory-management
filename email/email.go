package email

import (
	"github.com/Timothylock/inventory-management/config"

	"gopkg.in/gomail.v2"
)

type Sender interface {
	DialAndSend(m ...*gomail.Message) error
}

type Service struct {
	cfg    config.Config
	client Sender
}

func NewService(c config.Config, d Sender) Service {
	return Service{
		cfg:    c,
		client: d,
	}
}

func (s *Service) SendEmail(toEmail, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.EmailFromAddr)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.client.DialAndSend(m)
}
