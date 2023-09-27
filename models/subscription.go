package models

type Subscription struct {
	Model
	UserId               uint   `json:"user_id" gorm:"type:varchar(255);not null;index"`
	PodcastId            uint   `json:"podcast_id" gorm:"type:varchar(255);not null;index"`
	StripeSubscriptionId string `json:"-" gorm:"varchar(255);not null;unique"`
	Status               string `json:"status" gorm:"type:subscriptionstatus" sql:"type:SubscriptionStatus"`

	Listener User    `json:"listener" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Podcast  Podcast `json:"podcast" gorm:"foreignKey:PodcastId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
