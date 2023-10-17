package stripe

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/testhelpers/testclock"
)

const APPLICATION_FEE_PRECENT float64 = 20.00
const CREATOR_SHARE float64 = 100.00 - APPLICATION_FEE_PRECENT

type StripeGateway struct{}

type AccountLinkParams struct {
	ReturnUrl  string
	RefreshUrl string
}

type CheckoutSessionParams struct {
	UserId           string
	CustomerId       string
	CreatorAccountId string
	PodcastId        string
	SuccessUrl       string
	CancelUrl        string
}

type BillingSessionParams struct {
	CustomerId string
	ReturnUrl  string
}

func InitializeStripeGateway() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	log.Println(fmt.Sprintf("initialized the stripe gateway with version %s", stripe.ClientVersion))
}

func NewStripeGateway() *StripeGateway {
	return &StripeGateway{}
}

func createTestClock(uid string) *string {
	tc, _ := testclock.New(&stripe.TestHelpersTestClockParams{
		Name:       stripe.String("test_clock_" + uid),
		FrozenTime: stripe.Int64(time.Now().Unix()),
	})

	return &tc.ID
}
