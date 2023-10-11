package routes

import (
	"podcast/middleware"
	"podcast/types"

	"github.com/gin-gonic/gin"
)

func usersRoutes(r *gin.RouterGroup) {
	g := r.Group("/users")
	pc := CreatePodcastsController()

	g.Use(middleware.Authenticate())

	// creator group
	cg := g.Group("/")
	cg.Use(middleware.Authorize([]types.Role{types.CREATOR_ROLE}))
	cg.GET("/:id/podcast", pc.GetPodcastByCreator)

	// listener group
	lg := g.Group("/")
	lg.Use(middleware.Authorize([]types.Role{types.LISTENER_ROLE}))
	lg.GET("/:id/subscribed-podcasts", pc.GetListenerSubscribedPodcasts)
}
