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
	g.Use(middleware.Authorize([]types.Role{types.CREATOR_ROLE}))
	g.GET("/:cid/podcast", pc.GetPodcastByCreator)

}
