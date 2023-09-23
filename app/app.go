package app

import (
	"fmt"
	"os"
	"podcast/database"
	"podcast/routes"

	"github.com/gin-gonic/gin"
)

type App struct {
}

func init() {
	database.Connect()
	// database.Migrate() // should be enabled on only new database
}

func (a *App) Serve() {
	server := gin.Default()
	server.SetTrustedProxies(nil)

	g := server.Group("/api").Group("/v1")

	routes.RegisterRoutes(g)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server.Run(addr)
}
