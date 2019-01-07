package email

import (
	"errors"

	"github.com/Timothylock/inventory-management/config"

	"gopkg.in/gomail.v2"
)

type Service struct {
	client *gomail.Dialer
}

func NewService(c config.Config) Service {
	if c.EmailSMTPServ == "" {
		return Service{
			client: nil,
		}
	}

	return Service{
		client: gomail.NewPlainDialer(c.EmailSMTPServ, c.EmailSMTPPort, c.EmailUsername, c.EmailPassword),
	}
}

func (s *Service) sendEmail(toEmail, subject, body string) error {
	if s.client == nil {
		return errors.New("your admin did not set up email properly. Please contact them to reset your password for you")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.client.DialAndSend(m)
}
