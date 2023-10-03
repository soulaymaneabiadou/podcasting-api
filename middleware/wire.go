//go:build wireinject
// +build wireinject

package middleware

import (
	"podcast/database"
	"podcast/gateways/mailing"
	"podcast/repositories"
	"podcast/services"

	"github.com/google/wire"
)

func CreateAuthService() *services.AuthService {
	wire.Build(
		services.NewAuthService,
		wire.Bind(new(mailing.EmailGateway), new(*mailing.SMTPMailer)),
		services.NewEmailService,
		mailing.NewSMTPMailer,
		repositories.NewUsersRepository,
		database.Connection,
		// AuthServiceSet,
	)

	return &services.AuthService{}
}

// var UsersRepositorySet = wire.NewSet(
// 	repositories.NewUsersRepository,
// 	database.Connection,
// )

// var EmailServiceSet = wire.NewSet(
// 	wire.Bind(new(mailing.EmailGateway), new(*mailing.SMTPGateway)),
// 	services.NewEmailService,
// 	mailing.NewSMTPGateway,
// )

// var AuthServiceSet = wire.NewSet(
// 	services.NewAuthService,
// 	EmailServiceSet,
// 	UsersRepositorySet,
// )

// var AuthControllerSet = wire.newSet(
// 	controllers.NewAuthController
// 	AuthServiceSet,
// )

// var PodcastsRepositorySet = wire.NewSet(
// 	repositories.NewPodcastsRepository,
// 	database.Connection,
// )

// var PodcastsServiceSet = wire.NewSet(
// 	services.NewPodcastsService,
// 	PodcastsRepositorySet,
// )

// var PodcastsControllerSet = wire.NewSet(
// 	controllers.NewPodcastsController,
// 	PodcastsServiceSet,
// )
