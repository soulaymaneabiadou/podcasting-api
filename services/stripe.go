package services

import (
	"encoding/json"
	"fmt"
	"os"
	gateway "podcast/gateways/stripe"
	"podcast/repositories"
	"podcast/types"
	"strconv"

	"github.com/stripe/stripe-go/v75"
)

type StripeService struct {
	sg *gateway.StripeGateway
	sr *repositories.SubscriptionsRepository
	us *UsersService
}

type CustomerCheckoutSessionParams struct {
	UserId          string
	CustomerId      string
	StripeAccountId string
	PodcastId       string
}

func NewStripeService(sg *gateway.StripeGateway, sr *repositories.SubscriptionsRepository, us *UsersService) *StripeService {
	return &StripeService{sg: sg, sr: sr, us: us}
}

func (ss *StripeService) CreateAccount(u types.User) (*stripe.Account, error) {
	return ss.sg.CreateAccount(u)
}

func (ss *StripeService) GetAccount(aid string) (*stripe.Account, error) {
	return ss.sg.GetAccount(aid)
}

func (ss *StripeService) CreateAccountLink(acct *stripe.Account) (*stripe.AccountLink, error) {
	return ss.sg.CreateAccountLink(acct, gateway.AccountLinkParams{
		ReturnUrl:  os.Getenv("PUBLIC_URL") + "/api/v1/",
		RefreshUrl: os.Getenv("PUBLIC_URL") + "/api/v1/stripe/connect",
	})
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
		CancelUrl:        os.Getenv("PUBLIC_URL") + "/api/v1/podcasts/" + input.PodcastId + "?success=false",
	})

	return session.URL, err
}

func (ss *StripeService) CreateCustomerPortalSession(cid string) (*stripe.BillingPortalSession, error) {
	session, err := ss.sg.CreateCustomerPortalSession(gateway.BillingSessionParams{
		CustomerId: cid,
		ReturnUrl:  os.Getenv("PUBLIC_URL") + "/api/v1/auth/me",
	})

	return session, err
}

func (ss *StripeService) CreateConnectAccountLink(aid string) (*stripe.LoginLink, error) {
	session, err := ss.sg.CreateConnectAccountLink(aid)

	return session, err
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

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing subscription JSON: %v\n", err)
			return
		}

		ss.sr.UpdateStatus(subscription.ID, string(subscription.Status))

	case "account.updated":
		var account stripe.Account
		err := json.Unmarshal(event.Data.Raw, &account)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing account JSON: %v\n", err)
			return
		}

		ss.us.UpdateStripeInfo(account.ID, types.UpdateUserInput{
			ChargesEnabled:   account.ChargesEnabled,
			TransfersEnabled: account.PayoutsEnabled,
			DetailsSubmitted: account.DetailsSubmitted,
		})

	default:
		fmt.Fprintf(os.Stderr, "unhandled stripe event type: %s\n", event.Type)
	}
}
