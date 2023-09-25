package gateway

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SGMailer struct{}

func NewSGMailer() *SGMailer {
	return &SGMailer{}
}

func (s *SGMailer) SendEmail(p EmailPayload) {
	from := mail.NewEmail(os.Getenv("FROM_NAME"), os.Getenv("FROM_EMAIL"))

	to := mail.NewEmail(p.Receiver.Name, p.Receiver.Email)

	htmlContent := parseTemplate(p.Template, p.Data)

	message := mail.NewSingleEmail(from, p.Subject, to, htmlContent.String(), htmlContent.String())

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)
	if err != nil {
		log.Println("error occured while sending the email with sendgrid", err.Error())
	} else {
		fmt.Println("email sent via sendgrid", response)
	}
}
