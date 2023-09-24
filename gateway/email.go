package gateway

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"text/template"
)

type EmailPayload struct {
	Receiver string
	Subject  string
	Template string
	Data     any
}

func SendEmail(p EmailPayload) {
	identity := os.Getenv("SMTP_IDENTITY")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	to := []string{p.Receiver}
	body := parseTemplate(p.Template, p.Data)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte("From:" + from + "\r\n" +
		"To:" + p.Receiver + "\r\n" +
		"Subject:" + p.Subject + "\r\n" +
		mime +
		body.String())

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth(identity, username, password, host)

	err := smtp.SendMail(addr, auth, from, to, msg)
	if err == nil {
		log.Println("email has been sent")
	} else {
		log.Println("error occured while sending the email", err.Error())
	}
}

func parseTemplate(tpl string, data any) bytes.Buffer {
	t, err := parseTemplateDir("templates")
	if err != nil {
		log.Fatal("could not parse the templates dir", err)
	}

	var body bytes.Buffer
	err = t.ExecuteTemplate(&body, tpl, &data)
	if err != nil {
		log.Fatal("could not parse the template", err)
	}

	return body
}

func parseTemplateDir(dir string) (*template.Template, error) {
	var paths []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
