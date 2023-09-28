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
		wire.Bind(new(gateway.EmailGateway), new(*gateway.SMTPMailer)),
		gateway.NewSMTPMailer,
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

func CreateEpisodesController() *controllers.EpisodesController {
	wire.Build(
		controllers.NewEpisodesController,
		services.NewEpisodesService,
		repositories.NewEpisodesRepository,
		services.NewPodcastsService,
		repositories.NewPodcastsRepository,
		database.Connection,
		// EpisodesControllerSet,
	)

	return &controllers.EpisodesController{}
}

func CreateWebhooksController() *controllers.WebhooksController {
	wire.Build(
		controllers.NewWebhooksController,
		services.NewStripeService,
		gateway.NewStripeGateway,
		// WebhooksControllerSet,
	)

	return &controllers.WebhooksController{}
}
