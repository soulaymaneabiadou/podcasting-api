package repositories

import (
	"errors"
	"log"
	"podcast/types"

	"gorm.io/gorm"
)

type SubscriptionsRepository struct {
	db *gorm.DB
}

func NewSubscriptionsRepository(db *gorm.DB) *SubscriptionsRepository {
	return &SubscriptionsRepository{db: db}
}

func (sr *SubscriptionsRepository) GetAllByUserId(uid string) ([]types.Subscription, error) {
	var subscriptions []types.Subscription

	if err := sr.db.Where("user_id=?", uid).Find(&subscriptions).Error; err != nil {
		log.Println(err)
		return []types.Subscription{}, errors.New("no user subscriptions were found")
	}

	return subscriptions, nil
}

func (sr *SubscriptionsRepository) GetByUserAndPodcast(uid, pid string) (types.Subscription, error) {
	var subscription types.Subscription

	if err := sr.db.Where("user_id=?", uid).Where("podcast_id=?", pid).Find(&subscription).Error; err != nil {
		log.Println(err)
		return types.Subscription{}, errors.New("no user subscription was found for the given podcast")
	}

	return subscription, nil
}

func (sr *SubscriptionsRepository) Create(input types.CreateSubscriptionInput) (types.Subscription, error) {
	subscription := types.Subscription{
		UserId:               input.UserId,
		PodcastId:            input.PodcastId,
		StripeSubscriptionId: input.StripeSubscriptionId,
		Status:               input.Status,
	}

	if err := sr.db.Create(&subscription).Error; err != nil {
		log.Println(err)
		return types.Subscription{}, err
	}

	return subscription, nil
}

func (sr *SubscriptionsRepository) UpdateStatus(sid string, status string) (types.Subscription, error) {
	var subscription types.Subscription

	if err := sr.db.Model(&subscription).Where("stripe_subscription_id=?", sid).Update("status", status).Error; err != nil {
		log.Println(err)
		return types.Subscription{}, errors.New("subscription not found by stripe subscription id")
	}

	return subscription, nil
}
