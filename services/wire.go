//go:build wireinject
// +build wireinject

package services

import (
	"podcast/database"
	"podcast/gateways/mailing"
	"podcast/repositories"

	"github.com/google/wire"
)

func CreateAuthService() *AuthService {
	wire.Build(
		wire.Bind(new(mailing.EmailGateway), new(*mailing.SMTPMailer)),
		NewAuthService,
		NewEmailService,
		mailing.NewSMTPMailer,
		repositories.NewUsersRepository,
		database.Connection,
		// AuthServiceSet,
	)

	return &AuthService{}
}
