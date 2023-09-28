package services

import (
	"errors"
	"fmt"
	"strconv"

	"podcast/repositories"
	"podcast/types"
)

type PodcastsService struct {
	pr *repositories.PodcastsRepository
	us *UsersService
	ss *StripeService
}

func NewPodcastsService(
	pr *repositories.PodcastsRepository,
	us *UsersService,
	ss *StripeService,
) *PodcastsService {
	return &PodcastsService{pr: pr, us: us, ss: ss}
}

func (ps *PodcastsService) GetPodcasts(p types.Paginator) (int64, []types.Podcast, error) {
	count, err := ps.pr.Count()
	if err != nil {
		return 0, []types.Podcast{}, err
	}

	podcasts, err := ps.pr.GetAll(p)
	if err != nil {
		return 0, []types.Podcast{}, err
	}

	return count, podcasts, nil
}

func (ps *PodcastsService) GetPodcastById(id string) (types.Podcast, error) {
	podcast, err := ps.pr.GetById(id)
	if err != nil {
		return podcast, err
	}

	return podcast, nil
}

func (ps *PodcastsService) GetPodcastBySlug(uid string, slug string) (types.Podcast, error) {
	podcast, err := ps.pr.GetBySlug(slug)
	if err != nil {
		return podcast, err
	}

	return podcast, nil
}

func (ps *PodcastsService) CreatePodcast(d types.CreatePodcastInput) (types.Podcast, error) {
	return ps.pr.Create(d)
}

func (ps *PodcastsService) UpdatePodcast(uid string, id string, i types.UpdatePodcastInput) (types.Podcast, error) {
	podcast, err := ps.GetPodcastById(id)
	if err != nil {
		return podcast, err
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || podcast.CreatorId != uint(creatorId) {
		return podcast, errors.New("podcast not found")
	}

	podcast, err = ps.pr.Update(podcast, i)
	return podcast, err
}

func (ps *PodcastsService) DeletePodcast(uid string, id string) (bool, error) {
	podcast, err := ps.GetPodcastById(id)
	if err != nil {
		return false, err
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || podcast.CreatorId != uint(creatorId) {
		return false, errors.New("podcast not found")
	}

	return ps.pr.Destroy(podcast)
}

func (ps *PodcastsService) GetPodcastCreator(id string) (types.Account, error) {
	podcast, err := ps.GetPodcastById(id)
	if err != nil {
		return types.Account{}, err
	}

	account, err := ps.us.GetUserAccountById(fmt.Sprint(podcast.CreatorId))

	return account, nil
}

func (ps *PodcastsService) Subscribe(uid, pid string) (string, error) {
	user, err := ps.us.GetUserById(uid)
	if err != nil {
		return "", err
	}

	sub, err := ps.us.GetUserSubscriptionByPodcast(user, pid)
	if err != nil {
		return "", err
	}

	if sub.ID != 0 {
		return "", errors.New("you have already subscribed to this podcast")
	}

	if user.StripeCustomerId == "" {
		cus, err := ps.ss.CreateCustomer(user)
		if err != nil {
			return "", err
		}
		user, err = ps.us.SetUserCustomerId(user, cus.ID)
		if err != nil {
			return "", err
		}
	}

	account, err := ps.GetPodcastCreator(pid)
	if err != nil {
		return "", err
	}

	url, err := ps.ss.CreateCustomerCheckoutSession(CustomerCheckoutSessionParams{
		UserId:           fmt.Sprint(user.ID),
		CustomerId:       user.StripeCustomerId,
		CreatorAccountId: account.StripeAccountId,
		PodcastId:        pid,
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
