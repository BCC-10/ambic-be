package email

import (
	"ambic/internal/domain/env"
	"ambic/internal/infra/file"
	"fmt"
	gomail "gopkg.in/mail.v2"
)

type EmailIf interface {
	SendVerificationEmail(to string, code string) error
	SendResetPasswordLink(to string, token string) error
	SendPartnerRegistrationEmail(to string, name string) error
	SendPartnerVerifiedEmail(to string, name string) error
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

func (e *Email) SendVerificationEmail(to string, token string) error {
	message := gomail.NewMessage()
	body, err := file.ReadHTML(e.template, "verification")
	if err != nil {
		return err
	}

	message.SetHeader("From", e.user)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Verifikasi Email")

	message.SetBody("text/html", fmt.Sprintf(body, e.appURL, to, token))

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
	message.SetHeader("Subject", "Reset Password")

	message.SetBody("text/html", fmt.Sprintf(body, e.appURL, token))

	return e.sendEmail(e.connect(), message)
}

func (e *Email) SendPartnerRegistrationEmail(to string, name string) error {
	message := gomail.NewMessage()
	body, err := file.ReadHTML(e.template, "partner_registration")
	if err != nil {
		return err
	}

	message.SetHeader("From", e.user)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Pendaftaran Partner Berhasil")

	message.SetBody("text/html", fmt.Sprintf(body, name))

	return e.sendEmail(e.connect(), message)
}

func (e *Email) SendPartnerVerifiedEmail(to string, name string) error {
	message := gomail.NewMessage()
	body, err := file.ReadHTML(e.template, "partner_verified")
	if err != nil {
		return err
	}

	message.SetHeader("From", e.user)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Verifikasi Partner Berhasil")

	message.SetBody("text/html", fmt.Sprintf(body, name))

	return e.sendEmail(e.connect(), message)
}
