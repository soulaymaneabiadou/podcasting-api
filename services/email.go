package services

import (
	"fmt"
	"os"
	"podcast/gateways/mailing"
	"podcast/types"
	"time"
)

type EmailService struct {
	g mailing.EmailGateway
}

func NewEmailService(g mailing.EmailGateway) *EmailService {
	return &EmailService{g: g}
}

func (es *EmailService) createVerificationUrl(token string) string {
	return fmt.Sprintf("%s/auth/verify/%s", os.Getenv("PUBLIC_URL"), token)
}

func (es *EmailService) SendListenerWelcomeEmail(user types.User, token string) {
	subject := "Welcome to Podcasting Platform"
	data := map[string]string{
		"Subject":         subject,
		"Name":            user.Name,
		"VerificationURL": es.createVerificationUrl(token),
	}

	es.g.SendEmail(mailing.EmailPayload{
		Receiver: user,
		Subject:  subject,
		Template: "welcome.html",
		Data:     data,
	})
}

func (es *EmailService) SendCreatorWelcomeEmail(user types.User, token string) {
	subject := "Thanks for Joining Podcasting Platform"
	data := map[string]string{
		"Subject":         subject,
		"Name":            user.Name,
		"VerificationURL": es.createVerificationUrl(token),
	}

	es.g.SendEmail(mailing.EmailPayload{
		Receiver: user,
		Subject:  subject,
		Template: "welcome.html",
		Data:     data,
	})
}

func (es *EmailService) SendPasswordUpdatedEmail(user types.User) {
	subject := "[Podcasting Platform] Password Updated"
	data := map[string]string{
		"Subject": subject,
		"Name":    user.Name,
		"Date":    time.Now().Format("01-02-2006 15:04:05"),
		"URL":     os.Getenv("PUBLIC_URL") + "/auth/me",
	}

	es.g.SendEmail(mailing.EmailPayload{
		Receiver: user,
		Subject:  subject,
		Template: "password-updated.html",
		Data:     data,
	})
}

func (es *EmailService) SendPasswordResetTokenEmail(user types.User, token string) {
	subject := "[Podcasting Platform] Password Reset Link"
	data := map[string]string{
		"Subject": subject,
		"Name":    user.Name,
		"URL":     os.Getenv("PUBLIC_URL") + "/auth/resetpassword/" + token,
	}

	es.g.SendEmail(mailing.EmailPayload{
		Receiver: user,
		Subject:  subject,
		Template: "password-reset.html",
		Data:     data,
	})
}

func (es *EmailService) SendPasswordResettedEmail(user types.User) {
	subject := "[Podcasting Platform] Password Resetted Successfully"
	data := map[string]string{
		"Subject": subject,
		"Name":    user.Name,
		"URL":     os.Getenv("PUBLIC_URL") + "/auth/me",
	}

	es.g.SendEmail(mailing.EmailPayload{
		Receiver: user,
		Subject:  subject,
		Template: "password-resetted.html",
		Data:     data,
	})
}
