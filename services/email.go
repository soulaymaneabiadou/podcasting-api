package services

import (
	"os"
	"podcast/gateway"
	"podcast/types"
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
