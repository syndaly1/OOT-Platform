package usecase

import (
	"fmt"
	"net/smtp"
	"streaming-platform/auth-service/config"
)

type EmailSender struct {
	cfg *config.Config
}

func NewEmailSender(cfg *config.Config) *EmailSender {
	return &EmailSender{cfg: cfg}
}

func (e *EmailSender) SendVerificationEmail(to, token string) error {
	subject := "Verify your email"
	body := fmt.Sprintf("Click to verify: http://localhost:8080/auth/verify?token=%s", token)
	return e.send(to, subject, body)
}

func (e *EmailSender) SendPasswordResetEmail(to, token string) error {
	subject := "Reset your password"
	body := fmt.Sprintf("Click to reset: http://localhost:8080/auth/reset-password?token=%s", token)
	return e.send(to, subject, body)
}

func (e *EmailSender) send(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.cfg.SMTP.User, e.cfg.SMTP.Password, e.cfg.SMTP.Host)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")
	addr := fmt.Sprintf("%s:%s", e.cfg.SMTP.Host, e.cfg.SMTP.Port)
	return smtp.SendMail(addr, auth, e.cfg.SMTP.User, []string{to}, msg)
}
