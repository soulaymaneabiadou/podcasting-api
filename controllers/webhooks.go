package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"podcast/services"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75/webhook"
)

type WebhooksController struct {
	ss *services.StripeService
}

func NewWebhooksController(ss *services.StripeService) *WebhooksController {
	return &WebhooksController{ss: ss}
}

func (wc *WebhooksController) HandleStripeWebhooks(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading request body: %v\n", err)
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error verifying webhook signature: %v\n", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	wc.ss.HandleWebhookEvent(event)

	c.Writer.WriteHeader(http.StatusOK)
}
