package services

import (
	"errors"
	"podcast/repositories"
	"podcast/types"
	"strconv"
)

type PodcastsService struct {
	pr *repositories.PodcastsRepository
}

func NewPodcastsService(pr *repositories.PodcastsRepository) *PodcastsService {
	return &PodcastsService{pr: pr}
}

func (us *PodcastsService) GetPodcasts(p types.Paginator) (int64, []types.Podcast, error) {
	count, err := us.pr.Count()
	if err != nil {
		return 0, []types.Podcast{}, err
	}

	podcasts, err := us.pr.GetAll(p)
	if err != nil {
		return 0, []types.Podcast{}, err
	}

	return count, podcasts, nil
}

func (us *PodcastsService) GetPodcastById(id string) (types.Podcast, error) {
	podcast, err := us.pr.GetById(id)
	if err != nil {
		return podcast, err
	}

	return podcast, nil
}

func (us *PodcastsService) GetPodcastBySlug(uid string, slug string) (types.Podcast, error) {
	podcast, err := us.pr.GetBySlug(slug)
	if err != nil {
		return podcast, err
	}

	return podcast, nil
}

func (us *PodcastsService) CreatePodcast(d types.CreatePodcastInput) (types.Podcast, error) {
	return us.pr.Create(d)
}

func (us *PodcastsService) UpdatePodcast(uid string, id string, i types.UpdatePodcastInput) (types.Podcast, error) {
	podcast, err := us.GetPodcastById(id)
	if err != nil {
		return podcast, err
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || podcast.CreatorId != uint(creatorId) {
		return podcast, errors.New("podcast not found")
	}

	podcast, err = us.pr.Update(podcast, i)
	return podcast, err
}

func (us *PodcastsService) DeletePodcast(uid string, id string) (bool, error) {
	podcast, err := us.GetPodcastById(id)
	if err != nil {
		return false, err
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || podcast.CreatorId != uint(creatorId) {
		return false, errors.New("podcast not found")
	}

	return us.pr.Destroy(podcast)
}
