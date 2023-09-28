package routes

import (
	"podcast/controllers"

	"github.com/gin-gonic/gin"
)

func webhooksRoutes(r *gin.RouterGroup) {
	wc := controllers.NeWebhooksController()

	r.POST("/webhooks/stripe", wc.HandleStripeWebhooks)
}
