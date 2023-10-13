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

func (ps *PodcastsService) GetPodcasts(f types.PodcastFilters, s types.Sorter, p types.Paginator) (int64, []types.Podcast, error) {
	count, err := ps.pr.Count()
	if err != nil {
		return 0, []types.Podcast{}, err
	}

	podcasts, err := ps.pr.GetAll(f, s, p)
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

func (ps *PodcastsService) GetPodcastBySlug(slug string) (types.Podcast, error) {
	podcast, err := ps.pr.GetBySlug(slug)
	if err != nil {
		return podcast, err
	}

	return podcast, nil
}

func (ps *PodcastsService) CreatePodcast(d types.CreatePodcastInput) (types.Podcast, error) {
	podcast, err := ps.pr.Create(d)
	if err != nil {
		return types.Podcast{}, errors.New("you cannot create more than one podcast")
	}

	return podcast, nil
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

func (ps *PodcastsService) GetPodcastCreator(id string) (types.User, error) {
	podcast, err := ps.GetPodcastById(id)
	if err != nil {
		return types.User{}, err
	}

	return ps.us.GetUserById(fmt.Sprint(podcast.CreatorId))
}

func (ps *PodcastsService) GetPodcastByCreatorId(id string) (types.Podcast, error) {
	podcast, err := ps.pr.GetByCreatorId(id)
	if err != nil {
		return podcast, err
	}

	return podcast, nil
}

func (ps *PodcastsService) GetStats(id uint) (types.PodcastStats, error) {
	pid := fmt.Sprintf("%d", id)

	episodesCount, err := ps.pr.GetEpisodesCount(pid)
	if err != nil {
		return types.PodcastStats{}, err
	}

	activeSubCount, err := ps.pr.GetActiveSubscriptionsCount(pid)
	if err != nil {
		return types.PodcastStats{}, err
	}

	subCount, err := ps.pr.GetSubscriptionsCount(pid)
	if err != nil {
		return types.PodcastStats{}, err
	}

	// TODO:
	// stripe stats

	// construct stats
	stats := types.PodcastStats{
		TotalSubsCount:  subCount,
		ActiveSubsCount: activeSubCount,
		EpisodesCount:   episodesCount,
	}

	return stats, nil
}

func (ps *PodcastsService) GetListenerSubscribedPodcasts(id string) ([]types.Podcast, error) {
	podcasts, err := ps.pr.GetByListenerId(id)
	if err != nil {
		return []types.Podcast{}, err
	}

	return podcasts, nil
}
