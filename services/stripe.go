package services

import (
	"fmt"
	"log"
	"os"
	"podcast/gateway"

	"github.com/stripe/stripe-go/v75"
)

type StripeService struct {
	sg *gateway.StripeGateway
}

func NewStripeService(sg *gateway.StripeGateway) *StripeService {
	return &StripeService{sg: sg}
}

func (ss *StripeService) HandleWebhookEvent(event stripe.Event) {
	switch event.Type {
	case "customer.subscription.created":
		log.Println("sub created", event.Data.Object)

	default:
		fmt.Fprintf(os.Stderr, "unhandled stripe event type: %s\n", event.Type)
	}
}
