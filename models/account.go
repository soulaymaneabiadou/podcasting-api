package models

type Account struct {
	Model
	UserId           uint   `json:"user_id" gorm:"unique;not null"`
	StripeAccountId  string `json:"stripe_account_id" gorm:"type:varchar(255);not null;unique"`
	ChargesEnabled   bool   `json:"charges_enabled" gorm:"default:false"`
	TransfersEnabled bool   `json:"transfers_enabled" gorm:"default:false"`
	DetailsSubmitted bool   `json:"detailsSubmitted" gorm:"default:false"`
}
