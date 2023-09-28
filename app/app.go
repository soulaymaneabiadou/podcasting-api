package app

import (
	"fmt"
	"os"
	"podcast/database"
	"podcast/gateway"
	"podcast/routes"

	"github.com/gin-gonic/gin"
)

type App struct {
}

func init() {
	database.Connect()
	database.Migrate() // should be enabled only on a new database
	gateway.InitializeStripeGateway()
}

func (a *App) Serve() {
	server := gin.Default()
	server.SetTrustedProxies(nil)

	g := server.Group("/api").Group("/v1")

	routes.RegisterRoutes(g)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server.Run(addr)
}
