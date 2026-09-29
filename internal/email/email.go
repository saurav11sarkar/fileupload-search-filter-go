package email

import (
	"fmt"
	"net/smtp"

	"github.com/saurav11sarkar/001practic/internal/config"
)

type Email struct {
	cfg config.Config
}

func NewEmail(cfg config.Config) *Email {
	return &Email{cfg: cfg}
}

func (e *Email) SendEmail(to, subject, body string) error {
	if e.cfg.Email.Host == "" ||
		e.cfg.Email.Username == "" ||
		e.cfg.Email.Password == "" ||
		e.cfg.Email.From == "" {
		return fmt.Errorf("SMTP not configured")
	}

	auth := smtp.PlainAuth(
		"",
		e.cfg.Email.Username,
		e.cfg.Email.Password,
		e.cfg.Email.Host,
	)

	msg := fmt.Appendf(
		nil,
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		to,
		e.cfg.Email.From,
		subject,
		body,
	)

	addr := fmt.Sprintf("%s:%d", e.cfg.Email.Host, e.cfg.Email.Port)

	return smtp.SendMail(addr, auth, e.cfg.Email.From, []string{to}, msg)
}
