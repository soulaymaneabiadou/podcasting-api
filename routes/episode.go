package routes

import (
	"podcast/middleware"
	"podcast/types"

	"github.com/gin-gonic/gin"
)

func episodesRoutes(r *gin.RouterGroup) {
	g := r.Group("/episodes")
	ec := CreateEpisodesController()

	g.Use(middleware.Authenticate())

	g.GET("/:eid", ec.GetPodcastEpisode)

	g.Use(middleware.Authorize([]types.Role{types.CREATOR_ROLE}))
	g.POST("/:eid/publish", ec.PublishPodcastEpisode)
	g.PATCH("/:eid", ec.UpdatePodcastEpisode)
	g.DELETE("/:eid", ec.DeletePodcastEpisode)
}
