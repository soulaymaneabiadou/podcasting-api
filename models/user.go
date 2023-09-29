package models

import (
	"html"
	"strings"
	"time"

	"podcast/hasher"

	"gorm.io/gorm"
)

type User struct {
	Model
	Name                string    `json:"name" gorm:"type:varchar(255)"`
	Email               string    `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password            string    `json:"-" gorm:"type:varchar(255);not null"`
	ResetPasswordToken  string    `json:"-" gorm:"type:varchar(255);null;unique"`
	ResetPasswordExpire time.Time `json:"-"`
	Role                string    `json:"role" gorm:"type:role;default:'listener'" sql:"type:Role"`
	StripeCustomerId    string    `json:"stripe_customer_id" gorm:"type:varchar(255);null;unique"`
	StripeAccountId     string    `json:"stripe_account_id" gorm:"type:varchar(255);null;unique"`
	ChargesEnabled      bool      `json:"charges_enabled" gorm:"default:false"`
	TransfersEnabled    bool      `json:"transfers_enabled" gorm:"default:false"`
	DetailsSubmitted    bool      `json:"details_submitted" gorm:"default:false"`
	VerificationToken   string    `json:"-" gorm:"type:varchar(255);null;unique"`
	Verified            bool      `json:"-" gorm:"default:false"`
	VerifiedAt          time.Time `json:"-"`
	SigninCount         int64     `json:"-" gorm:"default:0;not null"`
	CurrentSigninAt     time.Time `json:"-"`
	CurrentSigninIP     string    `json:"-"`
	LastSigninAt        time.Time `json:"-"`
	LastSigninIP        string    `json:"-"`

	Subscriptions      []Subscription `json:"-"`
	SubscribedPodcasts []Podcast      `json:"subscribed_podcasts" gorm:"many2many:subscriptions"`
}

func (u *User) BeforeSave(db *gorm.DB) error {
	var err error

	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	u.Password, err = hasher.HashPassword(u.Password)

	return err
}
