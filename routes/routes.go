package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup) {
	healthRoutes(g)

	// public
	authRoutes(g)

	// auth middleware
	// g.Use(middleware.Authenticate())

	// private
	// userRoutes(g)
}
