//go:build wireinject
// +build wireinject

package routes

import (
	"podcast/controllers"
	"podcast/database"
	"podcast/gateways/mailing"
	"podcast/gateways/stripe"
	"podcast/gateways/upload"
	"podcast/repositories"
	"podcast/services"

	"github.com/google/wire"
)

func CreateAuthController() *controllers.AuthController {
	wire.Build(
		controllers.NewAuthController,
		services.NewAuthService,
		services.NewEmailService,
		wire.Bind(new(mailing.EmailGateway), new(*mailing.SMTPMailer)),
		mailing.NewSMTPMailer,
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
		wire.Bind(new(upload.FileHandler), new(*upload.LocalFileHandler)),
		upload.NewLocalFileHandler,
		// wire.Bind(new(upload.FileHandler), new(*upload.S3Handler)),
		// upload.NewS3Handler,
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

func CreateSubscriptionsController() *controllers.SubscriptionsController {
	wire.Build(
		controllers.NewSubscriptionsController,
		services.NewSubscriptionsService,
		podcastUserStripeSet,
		// WebhooksControllerSet,
	)

	return &controllers.SubscriptionsController{}
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
	stripe.NewStripeGateway,
	repositories.NewSubscriptionsRepository,
	database.Connection,
)
