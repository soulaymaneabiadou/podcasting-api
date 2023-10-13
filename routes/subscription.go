package routes

import (
	"podcast/middleware"
	"podcast/types"

	"github.com/gin-gonic/gin"
)

func subscriptionssRoutes(r *gin.RouterGroup) {
	g := r.Group("/subscriptions")
	sc := CreateSubscriptionsController()

	g.Use(middleware.Authenticate())

	g.POST("/", middleware.Authorize([]types.Role{types.LISTENER_ROLE}), sc.SubscribeToPodcast)
}
