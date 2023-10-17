package types

import (
	"podcast/models"
	"time"
)

type User = models.User

type Role string

const (
	CREATOR_ROLE  Role = "creator"
	LISTENER_ROLE Role = "listener"
)

type CreateUserInput struct {
	Name              string `json:"name" binding:"required,min=2"`
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required,min=8,alphanum"`
	Role              Role   `json:"role"`
	VerificationToken string `json:"-" binding:"-"`
}

type UpdateUserInput struct {
	Name                string    `json:"name" binding:"min=2"`
	Email               string    `json:"email" binding:"email"`
	Password            string    `json:"password" binding:"min=8,alphanum"`
	ResetPasswordToken  string    `json:"-"`
	ResetPasswordExpire time.Time `json:"-"`
	StripeCustomerId    string    `json:"stripe_customer_id" binding:"-" gorm:"unique"`
	StripeAccountId     string    `json:"stripe_account_id" binding:"-" gorm:"unique"`
	ChargesEnabled      bool      `json:"charges_enabled" binding:"-"`
	PayoutsEnabled      bool      `json:"transfers_enabled" binding:"-"`
	DetailsSubmitted    bool      `json:"details_submitted" binding:"-"`
	VerificationToken   string    `json:"-"`
	Verified            bool      `json:"-"`
	VerifiedAt          time.Time `json:"-"`
	SigninCount         uint      `json:"-"`
	CurrentSigninAt     time.Time `json:"-"`
	CurrentSigninIP     string    `json:"-"`
	LastSigninAt        time.Time `json:"-"`
	LastSigninIP        string    `json:"-"`
}
