package gateway

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type SMTPMailer struct{}

func NewSMTPMailer() *SMTPMailer {
	return &SMTPMailer{}
}

func (s *SMTPMailer) SendEmail(p EmailPayload) {
	identity := os.Getenv("SMTP_IDENTITY")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromName := os.Getenv("FROM_NAME")
	fromEmail := os.Getenv("FROM_EMAIL")

	to := []string{p.Receiver.Email}
	body := parseTemplate(p.Template, p.Data)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte("From:" + fromName + "<" + fromEmail + ">" + "\r\n" +
		"To:" + p.Receiver.Email + "\r\n" +
		"Subject:" + p.Subject + "\r\n" +
		mime +
		body.String())

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth(identity, username, password, host)

	err := smtp.SendMail(addr, auth, fromEmail, to, msg)
	if err == nil {
		log.Println("email has been sent via smtp")
	} else {
		log.Println("error occured while sending the email via smtp", err.Error())
	}
}
