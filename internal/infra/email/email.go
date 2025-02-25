package email

import (
	"ambic/internal/domain/env"
	"fmt"
	"net/smtp"
)

type EmailIf interface {
	SendOTP(to string, code string) error
}

type Email struct {
	host string
	port string
	user string
	pass string
}

type Config struct {
}

func NewEmail(env *env.Env) EmailIf {
	host := env.SMTPHost
	port := env.SMTPPort
	user := env.SMTPUser
	pass := env.SMTPPass

	return &Email{
		host, port, user, pass,
	}
}

func (e *Email) SendOTP(to string, otp string) error {
	subject := "Subject: Email Verification Code \n"
	body := fmt.Sprintf("Your verification code is %s", otp)
	message := []byte(subject + "\n" + body)

	return e.sendEmail(to, message)
}

func (e *Email) connect() smtp.Auth {
	return smtp.PlainAuth("", e.user, e.pass, e.host)
}

func (e *Email) sendEmail(to string, message []byte) error {
	return smtp.SendMail(e.host+":"+e.port, e.connect(), e.user, []string{to}, message)
}
