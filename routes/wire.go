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
	)

	return &controllers.PodcastsController{}
}

func CreateEpisodesController() *controllers.EpisodesController {
	wire.Build(
		controllers.NewEpisodesController,
		services.NewEpisodesService,
		repositories.NewEpisodesRepository,
		podcastUserStripeSet,
	)

	return &controllers.EpisodesController{}
}

func CreateWebhooksController() *controllers.WebhooksController {
	wire.Build(
		controllers.NewWebhooksController,
		stripeServiceSet,
		// WebhooksControllerSet,
	)

	return &controllers.WebhooksController{}
}

func CreateStripeController() *controllers.StripeController {
	wire.Build(
		controllers.NewStripeController,
		stripeServiceSet,
		// WebhooksControllerSet,
	)

	return &controllers.StripeController{}
}

var podcastUserStripeSet = wire.NewSet(
	services.NewPodcastsService,
	repositories.NewPodcastsRepository,
	stripeServiceSet,
)

var stripeServiceSet = wire.NewSet(
	services.NewUsersService,
	repositories.NewUsersRepository,
	services.NewStripeService,
	gateway.NewStripeGateway,
	repositories.NewSubscriptionsRepository,
	database.Connection,
)
