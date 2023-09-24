package services

import (
	"os"
	"podcast/gateway"
	"podcast/types"
	"time"
)

func SendWelcomeEmail(user types.User) {
	subject := "Welcome to Podcasting Platform"
	data := map[string]string{
		"Subject": subject,
		"Name":    user.Name,
		"URL":     os.Getenv("PUBLIC_URL"),
	}

	gateway.SendEmail(gateway.EmailPayload{
		Receiver: user.Email,
		Subject:  subject,
		Template: "welcome.html",
		Data:     data,
	})
}

func SendPasswordUpdatedEmail(user types.User) {
	subject := "[Podcasting Platform] Password Updated"
	data := map[string]string{
		"Subject": subject,
		"Name":    user.Name,
		"Date":    time.Now().Format("01-02-2006 15:04:05"),
		"URL":     os.Getenv("PUBLIC_URL") + "/auth/me",
	}

	gateway.SendEmail(gateway.EmailPayload{
		Receiver: user.Email,
		Subject:  subject,
		Template: "password-updated.html",
		Data:     data,
	})
}

func SendPasswordResetTokenEmail(user types.User, token string) {
	subject := "[Podcasting Platform] Password Reset Link"
	data := map[string]string{
		"Subject": subject,
		"Name":    user.Name,
		"URL":     os.Getenv("PUBLIC_URL") + "/auth/resetpassword/" + token,
	}

	gateway.SendEmail(gateway.EmailPayload{
		Receiver: user.Email,
		Subject:  subject,
		Template: "password-reset.html",
		Data:     data,
	})
}

func SendPasswordResettedEmail(user types.User) {
	subject := "[Podcasting Platform] Password Resetted Successfully"
	data := map[string]string{
		"Subject": subject,
		"Name":    user.Name,
		"URL":     os.Getenv("PUBLIC_URL") + "/auth/me",
	}

	gateway.SendEmail(gateway.EmailPayload{
		Receiver: user.Email,
		Subject:  subject,
		Template: "password-resetted.html",
		Data:     data,
	})
}
