package routes

import (
	"github.com/gin-gonic/gin"
)

func webhooksRoutes(r *gin.RouterGroup) {
	wc := CreateWebhooksController()

	r.POST("/webhooks/stripe", wc.HandleStripeWebhooks)
}
