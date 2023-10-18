package routes

import (
	"podcast/middleware"
	"podcast/types"

	"github.com/gin-gonic/gin"
)

func stripeRoutes(r *gin.RouterGroup) {
	g := r.Group("", middleware.Authenticate())
	sc := CreateStripeController()

	ag := g.Group("/account", middleware.Authorize([]types.Role{types.CREATOR_ROLE}))
	ag.POST("/connect", sc.CreateAccount)
	ag.POST("/onboard", sc.OnboardAccount)
	ag.POST("/login", sc.CreateAccountLogin)

	cg := g.Group("/customer", middleware.Authorize([]types.Role{types.LISTENER_ROLE}))
	cg.POST("/portal", sc.CreateCustomerPortal)
}
