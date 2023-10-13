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
	g.GET("/slugs/:pslug", pc.GetPodcastBySlug)
	g.GET("slugs/:pslug/episodes", ec.GetPodcastEpisodesBySlug)

	g.Use(middleware.Authenticate())

	g.Use(middleware.Authorize([]types.Role{types.CREATOR_ROLE}))

	g.POST("/", pc.CreatePodcast)
	g.GET("/:pid", pc.GetPodcast)
	g.PATCH("/:pid", pc.UpdatePodcast)
	g.DELETE("/:pid", pc.DeletePodcast)

	g.GET("/:pid/episodes", ec.GetPodcastEpisodes)
	g.POST("/:pid/episodes", ec.CreatePodcastEpisode)
}
