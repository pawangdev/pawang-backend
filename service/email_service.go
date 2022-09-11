package service

import (
	"pawang-backend/config"

	"gopkg.in/gomail.v2"
)

var CONFIG_SMTP_HOST = config.GetEnv("CONFIG_SMTP_HOST")
var CONFIG_SMTP_PORT = 587
var CONFIG_SENDER_NAME = config.GetEnv("CONFIG_SENDER_NAME")
var CONFIG_AUTH_EMAIL = config.GetEnv("CONFIG_AUTH_EMAIL")
var CONFIG_AUTH_PASSWORD = config.GetEnv("CONFIG_AUTH_PASSWORD")

type EmailService interface {
	SendEmail(emailTo string, subject string, body string) error
}

type emailService struct {
}

func NewEmailService() *emailService {
	return &emailService{}
}

func (service *emailService) SendEmail(emailTo string, subject string, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", emailTo)
	mailer.SetAddressHeader("Cc", emailTo, subject)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(CONFIG_SMTP_HOST, CONFIG_SMTP_PORT, CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}
