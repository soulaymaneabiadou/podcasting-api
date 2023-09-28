package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup) {
	healthRoutes(g)
	authRoutes(g)
	podcastsRoutes(g)
	episodesRoutes(g)
	webhooksRoutes(g)
}
