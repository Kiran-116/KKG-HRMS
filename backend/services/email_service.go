package services

import (
	"fmt"

	"hrms/config"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
}

type emailService struct {
	enabled bool
	from    string
	dialer  *gomail.Dialer
}

func NewEmailService() EmailService {
	cfg := config.AppConfig.SMTP
	if !cfg.Enabled {
		return &emailService{enabled: false}
	}

	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	return &emailService{
		enabled: true,
		from:    cfg.From,
		dialer:  dialer,
	}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	if !s.enabled {
		// Log that email would be sent in production
		fmt.Printf("Email would be sent to %s: %s\n", to, subject)
		return nil
	}

	message := gomail.NewMessage()
	message.SetHeader("From", s.from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	return s.dialer.DialAndSend(message)
}
