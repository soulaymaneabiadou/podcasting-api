package app

import (
	"fmt"
	"os"

	"podcast/database"
	"podcast/gateways/stripe"
	"podcast/middleware"
	"podcast/routes"

	"github.com/gin-gonic/gin"
)

type App struct {
}

func init() {
	database.Connect()
	stripe.InitializeStripeGateway()
}

func (a *App) Migrate() {
	database.Migrate()
}

func (a *App) Seed() {
	database.Drop()
	a.Migrate()
	database.Seed()
}

func (a *App) Serve() {
	server := gin.New()
	server.SetTrustedProxies(nil)

	server.Use(
		gin.Recovery(),
		middleware.Secure(),
		middleware.Cors(),
		middleware.Compression(),
		middleware.RequestId(),
		middleware.HttpLogger(),
		middleware.RateLimit(100, 1),
	)

	g := server.Group("/api").Group("/v1")

	routes.RegisterRoutes(g)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server.Run(addr)
}
