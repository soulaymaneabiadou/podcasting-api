//go:build wireinject
// +build wireinject

package services

import (
	"podcast/database"
	"podcast/gateway"
	"podcast/repositories"

	"github.com/google/wire"
)

func CreateAuthService() *AuthService {
	wire.Build(
		wire.Bind(new(gateway.EmailGateway), new(*gateway.SMTPMailer)),
		NewAuthService,
		NewEmailService,
		gateway.NewSMTPMailer,
		repositories.NewUsersRepository,
		database.Connection,
		// AuthServiceSet,
	)

	return &AuthService{}
}
