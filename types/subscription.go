package types

import "podcast/models"

type Subscription = models.Subscription

type CreateSubscriptionInput struct {
	UserId               uint   `json:"user_id" binding:"-"`
	PodcastId            uint   `json:"podcast_id" binding:"-"`
	StripeSubscriptionId string `json:"-" binding:"-"`
	Status               string `json:"status" binding:"-"`
}

type SubscribeInput struct {
	PodcastId string `json:"podcast_id" binding:"-"`
}
