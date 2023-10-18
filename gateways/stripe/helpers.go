package stripe

import (
	"time"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/testhelpers/testclock"
)

func createTestClock(uid string) *string {
	tc, _ := testclock.New(&stripe.TestHelpersTestClockParams{
		Name:       stripe.String("test_clock_" + uid),
		FrozenTime: stripe.Int64(time.Now().Unix()),
	})

	return &tc.ID
}
