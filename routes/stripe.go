package routes

import (
	"podcast/middleware"
	"podcast/types"

	"github.com/gin-gonic/gin"
)

func stripeRoutes(r *gin.RouterGroup) {
	g := r.Group("/stripe")
	sc := CreateStripeController()

	g.Use(middleware.Authenticate(), middleware.Authorize([]types.Role{types.CREATOR_ROLE}))

	g.POST("/connect", sc.Connect)
	g.POST("/onboard", sc.Onboard)
}
