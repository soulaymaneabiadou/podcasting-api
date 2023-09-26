package services

import (
	"errors"
	"fmt"
	"podcast/repositories"
	"podcast/types"
	"strconv"
	"time"
)

type EpisodesService struct {
	er *repositories.EpisodesRepository
	ps *PodcastsService
}

func NewEpisodesService(er *repositories.EpisodesRepository, ps *PodcastsService) *EpisodesService {
	return &EpisodesService{er: er, ps: ps}
}

func (us *EpisodesService) GetPodcastEpisodes(pid string, p types.Paginator) (int64, []types.Episode, error) {
	count, err := us.er.Count(pid)
	if err != nil {
		return 0, []types.Episode{}, err
	}

	episodes, err := us.er.GetAll(pid, p)
	if err != nil {
		return 0, []types.Episode{}, err
	}

	return count, episodes, nil
}

func (us *EpisodesService) GetPodcastEpisodeById(id string) (types.Episode, error) {
	episode, err := us.er.GetById(id)
	if err != nil {
		return episode, err
	}

	return episode, nil
}

func (us *EpisodesService) GetPodcastEpisodeBySlug(slug string) (types.Episode, error) {
	episode, err := us.er.GetBySlug(slug)
	if err != nil {
		return episode, err
	}

	return episode, nil
}

func (us *EpisodesService) CreatePodcastEpisode(d types.CreateEpisodeInput) (types.Episode, error) {
	podcast, err := us.ps.GetPodcastById(fmt.Sprint(d.PodcastId))
	if err != nil {
		return types.Episode{}, err
	}

	if podcast.CreatorId != d.CreatorId {
		return types.Episode{}, errors.New("unauthorized creator")
	}

	episode, err := us.er.Create(d)
	if err != nil {
		return types.Episode{}, err
	}

	return episode, nil
}

func (us *EpisodesService) UpdatePodcastEpisode(uid, id string, i types.UpdateEpisodeInput) (types.Episode, error) {
	episode, err := us.GetPodcastEpisodeById(id)
	if err != nil {
		return episode, err
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || episode.CreatorId != uint(creatorId) {
		return episode, errors.New("podcast episode not found")
	}

	podcast, err := us.ps.GetPodcastById(fmt.Sprint(episode.PodcastId))
	if err != nil {
		return episode, err
	}

	if podcast.CreatorId != episode.CreatorId {
		return episode, errors.New("unauthorized creator, cannot update")
	}

	episode, err = us.er.Update(episode, i)
	return episode, err
}

func (us *EpisodesService) PublishPodcastEpisode(uid, id string, i types.PublishEpisodeInput) (types.Episode, error) {
	episode, err := us.GetPodcastEpisodeById(id)
	if err != nil {
		return episode, err
	}

	podcast, err := us.ps.GetPodcastById(fmt.Sprint(episode.PodcastId))
	if err != nil {
		return episode, err
	}

	if podcast.CreatorId != episode.CreatorId {
		return episode, errors.New("unauthorized creator, cannot publish")
	}

	if !episode.PublishedAt.IsZero() && episode.Visibility != "draft" {
		return episode, errors.New("cannot publish an already published episode")
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || episode.CreatorId != uint(creatorId) {
		return episode, errors.New("podcast episode not found")
	}

	episode, err = us.er.Update(episode, types.UpdateEpisodeInput{
		Visibility:  i.Visibility,
		PublishedAt: time.Now(),
	})
	return episode, err
}

func (us *EpisodesService) DeletePodcastEpisode(uid, id string) (bool, error) {
	episode, err := us.GetPodcastEpisodeById(id)
	if err != nil {
		return false, err
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || episode.CreatorId != uint(creatorId) {
		return false, errors.New("episode not found")
	}

	podcast, err := us.ps.GetPodcastById(fmt.Sprint(episode.PodcastId))
	if err != nil {
		return false, err
	}

	if podcast.CreatorId != episode.CreatorId {
		return true, errors.New("unauthorized creator, cannot delete")
	}

	if episode.Visibility != "draft" {
		return true, errors.New("cannot delete a non draft episode")
	}

	return us.er.Destroy(episode)
}
