package stripe

import (
	"fmt"
	"os"
	"podcast/types"

	"github.com/stripe/stripe-go/v75"
	billingsession "github.com/stripe/stripe-go/v75/billingportal/session"
	"github.com/stripe/stripe-go/v75/checkout/session"
	"github.com/stripe/stripe-go/v75/customer"
)

func (sg *StripeGateway) CreateCustomer(u types.User) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:     &u.Name,
		Email:    &u.Email,
		Metadata: map[string]string{"user_id": fmt.Sprint(u.ID)},
	}

	if os.Getenv("GIN_MODE") != "release" {
		params.TestClock = createTestClock(fmt.Sprint(u.ID))
	}

	c, err := customer.New(params)

	return c, err
}

func (sg *StripeGateway) CreateCustomerPortalSession(input BillingSessionParams) (*stripe.BillingPortalSession, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(input.CustomerId),
		ReturnURL: stripe.String(input.ReturnUrl),
	}

	s, err := billingsession.New(params)

	return s, err
}

func (sg *StripeGateway) CreateCustomerCheckoutSession(sp CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Mode:              stripe.String("subscription"),
		Currency:          stripe.String("usd"),
		ClientReferenceID: &sp.CustomerId,
		Customer:          &sp.CustomerId,
		CustomerUpdate: &stripe.CheckoutSessionCustomerUpdateParams{
			Address: stripe.String("auto"),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(os.Getenv("STRIPE_PRICE_ID")),
				Quantity: stripe.Int64(1),
			},
		},
		AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
			Enabled: stripe.Bool(true),
		},
		SuccessURL: stripe.String(sp.SuccessUrl),
		CancelURL:  stripe.String(sp.CancelUrl),
		Metadata: map[string]string{
			"user_id":    sp.UserId,
			"podcast_id": sp.PodcastId,
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id":    sp.UserId,
				"podcast_id": sp.PodcastId,
			},
			TransferData: &stripe.CheckoutSessionSubscriptionDataTransferDataParams{
				Destination:   &sp.CreatorAccountId,
				AmountPercent: stripe.Float64(CREATOR_SHARE),
			},
		},
	}

	s, err := session.New(params)

	return s, err
}
