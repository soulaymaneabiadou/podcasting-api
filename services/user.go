package services

import (
	"fmt"

	"podcast/repositories"
	"podcast/types"
)

type UsersService struct {
	ur *repositories.UsersRepository
	sr *repositories.SubscriptionsRepository
}

func NewUsersService(
	ur *repositories.UsersRepository,
	sr *repositories.SubscriptionsRepository,
) *UsersService {
	return &UsersService{ur: ur, sr: sr}
}

func (us *UsersService) GetUserById(id string) (types.User, error) {
	return us.ur.GetById(id)
}

func (us *UsersService) GetUserSubscriptionByPodcast(user types.User, pid string) (types.Subscription, error) {
	subscription, err := us.sr.GetByUserAndPodcast(fmt.Sprint(user.ID), pid)

	return subscription, err
}

func (us *UsersService) SetUserCustomerId(user types.User, cid string) (types.User, error) {
	return us.ur.Update(user, types.UpdateUserInput{StripeCustomerId: cid})
}

func (us *UsersService) Update(user types.User, input types.UpdateUserInput) (types.User, error) {
	return us.ur.Update(user, input)
}

func (us *UsersService) UpdateStripeInfo(aid string, input types.UpdateUserInput) (types.User, error) {
	user, err := us.ur.GetByStripeAccountId(aid)
	if err != nil {
		return user, err
	}

	return us.ur.Update(user, input)
}

func (us *UsersService) SetStripeAccountId(user types.User, aid string) (types.User, error) {
	return us.Update(user, types.UpdateUserInput{
		StripeAccountId:  aid,
		ChargesEnabled:   false,
		TransfersEnabled: false,
		DetailsSubmitted: false,
	})
}
