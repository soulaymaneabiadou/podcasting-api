package app

import (
	"fmt"
	"os"
	"podcast/database"
	"podcast/gateways/stripe"
	"podcast/routes"
	"time"

	"github.com/gin-contrib/cors"
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
	server := gin.Default()
	server.SetTrustedProxies(nil)

	config := cors.Config{
		AllowOrigins:     []string{os.Getenv("PUBLIC_URL")},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		MaxAge:           12 * time.Hour,
	}
	server.Use(cors.New(config))

	g := server.Group("/api").Group("/v1")

	routes.RegisterRoutes(g)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server.Run(addr)
}
