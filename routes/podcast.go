package routes

import (
	"podcast/controllers"
	"podcast/middleware"
	"podcast/types"

	"github.com/gin-gonic/gin"
)

func podcastsRoutes(r *gin.RouterGroup) {
	g := r.Group("/podcasts")
	// TODO: DI
	pc := controllers.NewPodcastsController()

	g.GET("/", pc.GetPodcasts)
	g.GET("/:id", pc.GetPodcast)

	g.Use(middleware.Authenticate(), middleware.Authorize([]types.Role{types.CREATOR_ROLE}))

	g.POST("/", pc.CreatePodcast)
	g.PATCH("/:id", pc.UpdatePodcast)
	g.DELETE("/:id", pc.DeletePodcast)
}
