//go:build wireinject
// +build wireinject

package routes

import (
	"podcast/controllers"
	"podcast/database"
	"podcast/gateway"
	"podcast/repositories"
	"podcast/services"

	"github.com/google/wire"
)

func CreateAuthController() *controllers.AuthController {
	wire.Build(
		controllers.NewAuthController,
		services.NewAuthService,
		services.NewEmailService,
		wire.Bind(new(gateway.EmailGateway), new(*gateway.SMTPGateway)),
		gateway.NewSMTPGateway,
		repositories.NewUsersRepository,
		database.Connection,
		// AuthControllerSet,
	)

	return &controllers.AuthController{}
}

func CreatePodcastsController() *controllers.PodcastsController {
	wire.Build(
		controllers.NewPodcastsController,
		services.NewPodcastsService,
		repositories.NewPodcastsRepository,
		database.Connection,
		// PodcastsControllerSet,
	)

	return &controllers.PodcastsController{}
}
