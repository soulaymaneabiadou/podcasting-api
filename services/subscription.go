package services

import (
	"errors"
	"fmt"
	"podcast/repositories"
)

type SubscriptionsService struct {
	pr *repositories.SubscriptionsRepository
	us *UsersService
	ss *StripeService
	ps *PodcastsService
}

func NewSubscriptionsService(
	pr *repositories.SubscriptionsRepository,
	us *UsersService,
	ss *StripeService,
	ps *PodcastsService,
) *SubscriptionsService {
	return &SubscriptionsService{pr: pr, us: us, ss: ss, ps: ps}
}

func (ss *SubscriptionsService) Subscribe(uid, pid string) (string, error) {
	podcast, err := ss.ps.GetPodcastById(pid)
	if err != nil {
		return "", err
	}

	user, err := ss.us.GetUserById(uid)
	if err != nil {
		return "", err
	}

	if user.Role != "listener" {
		return "", errors.New("unauthorized to subscribe, not a listener")
	}

	sub, err := ss.us.GetUserSubscriptionByPodcast(user, fmt.Sprint(podcast.ID))
	if err != nil {
		return "", err
	}

	if sub.ID != 0 {
		return "", errors.New("you have already subscribed to this podcast")
	}

	if user.StripeCustomerId == "" {
		cus, err := ss.ss.CreateCustomer(user)
		if err != nil {
			return "", err
		}
		user, err = ss.us.SetUserCustomerId(user, cus.ID)
		if err != nil {
			return "", err
		}
	}

	creator, err := ss.ps.GetPodcastCreator(fmt.Sprint(podcast.ID))
	if err != nil {
		return "", err
	}

	url, err := ss.ss.CreateCustomerCheckoutSession(CustomerCheckoutSessionParams{
		UserId:          fmt.Sprint(user.ID),
		PodcastId:       fmt.Sprint(podcast.ID),
		CustomerId:      user.StripeCustomerId,
		StripeAccountId: creator.StripeAccountId,
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
