package email

import (
	"ambic/internal/domain/env"
	"fmt"
	"net/smtp"
)

type EmailIf interface {
	SendOTP(to string, code string) error
	SendResetPassword(to string, token string) error
}

type Email struct {
	host   string
	port   string
	user   string
	pass   string
	appURL string
}

type Config struct {
}

func NewEmail(env *env.Env) EmailIf {
	host := env.SMTPHost
	port := env.SMTPPort
	user := env.SMTPUser
	pass := env.SMTPPass
	appURL := env.AppURL

	return &Email{
		host, port, user, pass, appURL,
	}
}

func (e *Email) connect() smtp.Auth {
	return smtp.PlainAuth("", e.user, e.pass, e.host)
}

func (e *Email) sendEmail(to string, message []byte) error {
	return smtp.SendMail(e.host+":"+e.port, e.connect(), e.user, []string{to}, message)
}

func (e *Email) SendOTP(to string, otp string) error {
	subject := "Subject: Email Verification Code \n"
	body := fmt.Sprintf("Your verification code is %s", otp)
	message := []byte(subject + "\n" + body)

	return e.sendEmail(to, message)
}

func (e *Email) SendResetPassword(to string, token string) error {
	subject := "Subject: Reset Password \n"
	body := fmt.Sprintf("Click this link to reset your password: %s/reset-password?token=%s", e.appURL, token)
	message := []byte(subject + "\n" + body)

	return e.sendEmail(to, message)
}
