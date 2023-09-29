package routes

import (
	"podcast/middleware"
	"podcast/types"

	"github.com/gin-gonic/gin"
)

func stripeRoutes(r *gin.RouterGroup) {
	g := r.Group("/stripe")
	sc := CreateStripeController()

	g.Use(middleware.Authenticate())

	g.POST("/connect", middleware.Authorize([]types.Role{types.CREATOR_ROLE}), sc.Connect)
	g.POST("/onboard", middleware.Authorize([]types.Role{types.CREATOR_ROLE}), sc.Onboard)

	g.POST("/portal", middleware.Authorize([]types.Role{types.LISTENER_ROLE}), sc.CustomerPortal)
}
