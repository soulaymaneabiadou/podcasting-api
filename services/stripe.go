package services

import (
	"encoding/json"
	"fmt"
	"os"
	"podcast/gateway"
	"podcast/repositories"
	"podcast/types"
	"strconv"

	"github.com/stripe/stripe-go/v75"
)

type StripeService struct {
	sg *gateway.StripeGateway
	sr *repositories.SubscriptionsRepository
}

type CustomerCheckoutSessionParams struct {
	UserId          string
	CustomerId      string
	StripeAccountId string
	PodcastId       string
}

func NewStripeService(sg *gateway.StripeGateway, sr *repositories.SubscriptionsRepository) *StripeService {
	return &StripeService{sg: sg, sr: sr}
}

func (ss *StripeService) CreateCustomer(u types.User) (*stripe.Customer, error) {
	return ss.sg.CreateCustomer(u)
}

func (ss *StripeService) CreateCustomerCheckoutSession(input CustomerCheckoutSessionParams) (string, error) {
	session, err := ss.sg.CreateCheckoutSession(gateway.CheckoutSessionParams{
		UserId:           input.UserId,
		CustomerId:       input.CustomerId,
		CreatorAccountId: input.StripeAccountId,
		PodcastId:        input.PodcastId,
		SuccessUrl:       os.Getenv("PUBLIC_URL") + "/api/v1/podcasts/" + input.PodcastId,
		CancelUrl:        os.Getenv("PUBLIC_URL") + "/api/v1/podcasts/" + input.PodcastId,
	})

	return session.URL, err
}

func (ss *StripeService) HandleWebhookEvent(event stripe.Event) {
	switch event.Type {
	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing subscription JSON: %v\n", err)
			return
		}

		uid, err := strconv.Atoi(subscription.Metadata["user_id"])
		pid, err := strconv.Atoi(subscription.Metadata["podcast_id"])

		ss.sr.Create(types.CreateSubscriptionInput{
			UserId:               uint(uid),
			PodcastId:            uint(pid),
			StripeSubscriptionId: subscription.ID,
			Status:               string(subscription.Status),
		})

	default:
		fmt.Fprintf(os.Stderr, "unhandled stripe event type: %s\n", event.Type)
	}
}
