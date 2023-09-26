package routes

import (
	"podcast/middleware"
	"podcast/types"

	"github.com/gin-gonic/gin"
)

func podcastsRoutes(r *gin.RouterGroup) {
	g := r.Group("/podcasts")
	pc := CreatePodcastsController()
	ec := CreateEpisodesController()

	g.GET("/", pc.GetPodcasts)
	g.GET("/:pid", pc.GetPodcast)
	g.GET("/:pid/episodes", ec.GetPodcastEpisodes)

	g.Use(middleware.Authenticate(), middleware.Authorize([]types.Role{types.CREATOR_ROLE}))

	g.POST("/", pc.CreatePodcast)
	g.PATCH("/:pid", pc.UpdatePodcast)
	g.DELETE("/:pid", pc.DeletePodcast)

	g.POST("/:pid/episodes", ec.CreatePodcastEpisode)
}
