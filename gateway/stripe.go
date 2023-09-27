package gateway

import (
	"fmt"
	"log"
	"os"
	"time"

	"podcast/types"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/account"
	"github.com/stripe/stripe-go/v75/accountlink"
	"github.com/stripe/stripe-go/v75/checkout/session"
	"github.com/stripe/stripe-go/v75/customer"
	"github.com/stripe/stripe-go/v75/testhelpers/testclock"
)

const APPLICATION_FEE_PRECENT float64 = 20.00
const CREATOR_SHARE float64 = 100.00 - APPLICATION_FEE_PRECENT

type StripeGateway struct{}

type AccountLinkParams struct {
	returnUrl  string
	refreshUrl string
}

type CheckoutSessionParams struct {
	customerId       string
	creatorAccountId string
	podcastId        string
	successUrl       string
	cancelUrl        string
}

func InitializeStripeGateway() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	log.Println(fmt.Sprintf("initialized the stripe gateway with version %s", stripe.ClientVersion))
}

func NewStripeGateway() *StripeGateway {
	return &StripeGateway{}
}

func (sg *StripeGateway) CreateAccount(u types.User) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Type:         stripe.String("express"),
		Country:      stripe.String("US"),
		Email:        stripe.String(u.Email),
		Metadata:     map[string]string{"user_id": fmt.Sprint(u.ID)},
		BusinessType: stripe.String("individual"),
		Individual:   &stripe.PersonParams{Email: &u.Email},
		BusinessProfile: &stripe.AccountBusinessProfileParams{
			MCC:                stripe.String("5815"),
			Name:               stripe.String(u.Name),
			SupportEmail:       stripe.String(u.Email),
			ProductDescription: stripe.String("Podcast creation"),
		},
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
		TOSAcceptance: &stripe.AccountTOSAcceptanceParams{
			ServiceAgreement: stripe.String("full"),
		},
	}

	a, err := account.New(params)

	return a, err
}

func (sg *StripeGateway) CreateAccountLink(acct *stripe.Account, p AccountLinkParams) (*stripe.AccountLink, error) {
	link, err := accountlink.New(&stripe.AccountLinkParams{
		Account:    &acct.ID,
		RefreshURL: stripe.String(p.refreshUrl),
		ReturnURL:  stripe.String(p.returnUrl),
		Type:       stripe.String("account_onboarding"),
		Collect:    stripe.String("eventually_due"),
	})

	return link, err
}

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

func (sg *StripeGateway) CreateCheckoutSession(sp CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Mode:              stripe.String("subscription"),
		Currency:          stripe.String("usd"),
		ClientReferenceID: &sp.customerId,
		Customer:          &sp.customerId,
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
		SuccessURL: stripe.String(sp.successUrl),
		CancelURL:  stripe.String(sp.cancelUrl),
		Metadata: map[string]string{
			"user_id":    sp.customerId,
			"creator_id": sp.creatorAccountId,
			"podcast_id": sp.podcastId,
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id":    sp.customerId,
				"creator_id": sp.creatorAccountId,
				"podcast_id": sp.podcastId,
			},
			OnBehalfOf: &sp.creatorAccountId,
			TransferData: &stripe.CheckoutSessionSubscriptionDataTransferDataParams{
				Destination:   &sp.creatorAccountId,
				AmountPercent: stripe.Float64(CREATOR_SHARE),
			},
		},
	}

	s, err := session.New(params)

	return s, err
}

func createTestClock(uid string) *string {
	tc, _ := testclock.New(&stripe.TestHelpersTestClockParams{
		Name:       stripe.String("test_clock_" + uid),
		FrozenTime: stripe.Int64(time.Now().Unix()),
	})

	return &tc.ID
}
