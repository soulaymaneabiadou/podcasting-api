package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75/webhook"
)

type WebhooksController struct {
}

func NeWebhooksController() *WebhooksController {
	return &WebhooksController{}
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

	switch event.Type {
	case "customer.subscription.created":
		log.Println("sub created", event.Data.Object)

	default:
		fmt.Fprintf(os.Stderr, "unhandled stripe event type: %s\n", event.Type)
	}

	c.Writer.WriteHeader(http.StatusOK)
}
