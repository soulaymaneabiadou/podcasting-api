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
		podcastUserStripeSet,
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
		podcastUserStripeSet,
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
		repositories.NewSubscriptionsRepository,
		database.Connection,
		// WebhooksControllerSet,
	)

	return &controllers.WebhooksController{}
}

var podcastUserStripeSet = wire.NewSet(
	services.NewPodcastsService,
	services.NewUsersService,
	services.NewStripeService,
	gateway.NewStripeGateway,
	repositories.NewUsersRepository,
	repositories.NewPodcastsRepository,
	repositories.NewSubscriptionsRepository,
	repositories.NewAccountsRepository,
)
