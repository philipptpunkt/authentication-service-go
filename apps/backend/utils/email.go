package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Dialer *gomail.Dialer
	From   string
}

func CreateEmailSender() (*EmailSender, error) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || user == "" || pass == "" || from == "" {
		return nil, fmt.Errorf("SMTP configuration is incomplete")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid smpt port: %v", err)
	}

	dialer := gomail.NewDialer(host, portInt, user, pass)
	return &EmailSender{Dialer: dialer, From: from}, nil
}

func (es *EmailSender) SendEmail(to, subject, body string, isHTML bool) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	if err := es.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
