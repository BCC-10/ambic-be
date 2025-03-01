package email

import (
	"ambic/internal/domain/env"
	"ambic/internal/infra/file"
	"fmt"
	gomail "gopkg.in/mail.v2"
)

type EmailIf interface {
	SendEmailVerification(to string, code string) error
	SendResetPasswordLink(to string, token string) error
}

type Email struct {
	host     string
	port     int
	user     string
	pass     string
	appURL   string
	template string
}

type Config struct {
}

func NewEmail(env *env.Env) EmailIf {
	host := env.SMTPHost
	port := env.SMTPPort
	user := env.SMTPUser
	pass := env.SMTPPass
	appURL := env.AppURL
	template := "internal/infra/email/template"

	return &Email{
		host, port, user, pass, appURL, template,
	}
}

func (e *Email) connect() *gomail.Dialer {
	return gomail.NewDialer(e.host, e.port, e.user, e.pass)
}

func (e *Email) sendEmail(dialer *gomail.Dialer, message *gomail.Message) error {
	return dialer.DialAndSend(message)
}

func (e *Email) SendEmailVerification(to string, token string) error {
	message := gomail.NewMessage()
	body, err := file.ReadHTML(e.template, "verification")
	if err != nil {
		return err
	}

	message.SetHeader("From", e.user)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Email Verification")

	message.SetBody("text/html", fmt.Sprintf(body, e.appURL, token))

	return e.sendEmail(e.connect(), message)
}

func (e *Email) SendResetPasswordLink(to string, token string) error {
	message := gomail.NewMessage()
	body, err := file.ReadHTML(e.template, "reset_password")
	if err != nil {
		return err
	}

	message.SetHeader("From", e.user)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Email Verification Code")

	message.SetBody("text/html", fmt.Sprintf(body, e.appURL, token))

	return e.sendEmail(e.connect(), message)
}
