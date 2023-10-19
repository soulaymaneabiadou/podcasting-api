package stripe

import (
	"fmt"
	"podcast/types"
	"strings"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/account"
	"github.com/stripe/stripe-go/v75/accountlink"
	"github.com/stripe/stripe-go/v75/balance"
	"github.com/stripe/stripe-go/v75/loginlink"
	"github.com/stripe/stripe-go/v75/payout"
)

func (sg *StripeGateway) CreateAccount(u types.User) (*stripe.Account, error) {
	nameParts := strings.Split(u.Name, " ")

	params := &stripe.AccountParams{
		Type:         stripe.String("express"),
		Country:      stripe.String("US"),
		Email:        stripe.String(u.Email),
		Metadata:     map[string]string{"user_id": fmt.Sprint(u.ID)},
		BusinessType: stripe.String("individual"),
		Individual: &stripe.PersonParams{
			Email:     stripe.String(u.Email),
			FirstName: stripe.String(nameParts[0]),
			LastName:  stripe.String(nameParts[1]),
		},
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
		Settings: &stripe.AccountSettingsParams{
			Payouts: &stripe.AccountSettingsPayoutsParams{
				Schedule: &stripe.AccountSettingsPayoutsScheduleParams{
					Interval:      stripe.String("monthly"),
					MonthlyAnchor: stripe.Int64(31),
				},
			},
		},
		TOSAcceptance: &stripe.AccountTOSAcceptanceParams{
			ServiceAgreement: stripe.String("full"),
		},
	}

	a, err := account.New(params)

	return a, err
}

func (sg *StripeGateway) GetAccount(id string) (*stripe.Account, error) {
	acct, err := account.GetByID(id, &stripe.AccountParams{})

	return acct, err
}

func (sg *StripeGateway) CreateAccountOnboardingLink(acct *stripe.Account, p AccountLinkParams) (*stripe.AccountLink, error) {
	link, err := accountlink.New(&stripe.AccountLinkParams{
		Account:    &acct.ID,
		RefreshURL: stripe.String(p.RefreshUrl),
		ReturnURL:  stripe.String(p.ReturnUrl),
		Type:       stripe.String("account_onboarding"),
		Collect:    stripe.String("eventually_due"),
	})

	return link, err
}

func (sg *StripeGateway) CreateAccountLoginLink(account string) (*stripe.LoginLink, error) {
	params := &stripe.LoginLinkParams{
		Account: stripe.String(account),
	}

	l, err := loginlink.New(params)

	return l, err
}

func (sg *StripeGateway) GetAccountBalance(accountId string) (*stripe.Balance, error) {
	params := &stripe.BalanceParams{}
	params.SetStripeAccount(accountId)

	result, err := balance.Get(params)

	return result, err
}

func (sg *StripeGateway) CreatePayout(accountId string, amount int64) (*stripe.Payout, error) {
	p, err := payout.New(&stripe.PayoutParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String("usd"),
		Method:   stripe.String("instant"),
		Params: stripe.Params{
			StripeAccount: &accountId,
		},
		// Destination: nil,
	})

	return p, err
}
