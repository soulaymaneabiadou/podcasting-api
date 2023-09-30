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
	us *UsersService
}

func NewEpisodesService(er *repositories.EpisodesRepository, ps *PodcastsService, us *UsersService) *EpisodesService {
	return &EpisodesService{er: er, ps: ps, us: us}
}

func (es *EpisodesService) GetPodcastEpisodes(pid string, p types.Paginator) (int64, []types.Episode, error) {
	count, err := es.er.Count(pid)
	if err != nil {
		return 0, []types.Episode{}, err
	}

	episodes, err := es.er.GetAll(pid, p)
	if err != nil {
		return 0, []types.Episode{}, err
	}

	return count, episodes, nil
}

func (es *EpisodesService) GetPodcastEpisodesBySlug(slug string, p types.Paginator) (int64, []types.Episode, error) {
	podcast, err := es.ps.GetPodcastBySlug(slug)
	if err != nil {
		return 0, []types.Episode{}, err
	}

	count, err := es.er.Count(fmt.Sprint(podcast.ID))
	if err != nil {
		return 0, []types.Episode{}, err
	}

	episodes, err := es.er.GetPublicEpisodes(fmt.Sprint(podcast.ID), p)
	if err != nil {
		return 0, []types.Episode{}, err
	}

	return count, episodes, nil
}

func (es *EpisodesService) GetPodcastEpisodeById(id string) (types.Episode, error) {
	episode, err := es.er.GetById(id)
	if err != nil {
		return episode, err
	}

	return episode, nil
}

func (es *EpisodesService) GetPodcastEpisodeBySlug(slug string) (types.Episode, error) {
	episode, err := es.er.GetPublicEpisodeBySlug(slug)
	if err != nil {
		return episode, err
	}

	return episode, nil
}

func (es *EpisodesService) CreatePodcastEpisode(d types.CreateEpisodeInput) (types.Episode, error) {
	podcast, err := es.ps.GetPodcastById(fmt.Sprint(d.PodcastId))
	if err != nil {
		return types.Episode{}, err
	}

	if podcast.CreatorId != d.CreatorId {
		return types.Episode{}, errors.New("unauthorized creator")
	}

	// check if creator has a stripe account if visibility is set to protected
	if d.Visibility == "protected" {
		creator, err := es.us.GetUserById(fmt.Sprint(podcast.CreatorId))
		if err != nil {
			return types.Episode{}, errors.New("an error occured while checking validity, please try again later")
		}

		if creator.StripeAccountId == "" || creator.DetailsSubmitted == false {
			return types.Episode{}, errors.New("you do not have a stripe account yet, please connect one before retrying thi action again.")
		}
	}

	episode, err := es.er.Create(d)
	if err != nil {
		return types.Episode{}, err
	}

	return episode, nil
}

func (es *EpisodesService) UpdatePodcastEpisode(uid, id string, i types.UpdateEpisodeInput) (types.Episode, error) {
	episode, err := es.GetPodcastEpisodeById(id)
	if err != nil {
		return episode, err
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || episode.CreatorId != uint(creatorId) {
		return episode, errors.New("podcast episode not found")
	}

	podcast, err := es.ps.GetPodcastById(fmt.Sprint(episode.PodcastId))
	if err != nil {
		return episode, err
	}

	if podcast.CreatorId != episode.CreatorId {
		return episode, errors.New("unauthorized creator, cannot update")
	}

	episode, err = es.er.Update(episode, i)
	return episode, err
}

func (es *EpisodesService) PublishPodcastEpisode(uid, id string, i types.PublishEpisodeInput) (types.Episode, error) {
	episode, err := es.GetPodcastEpisodeById(id)
	if err != nil {
		return episode, err
	}

	podcast, err := es.ps.GetPodcastById(fmt.Sprint(episode.PodcastId))
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

	episode, err = es.er.Update(episode, types.UpdateEpisodeInput{
		Visibility:  i.Visibility,
		PublishedAt: time.Now(),
	})
	return episode, err
}

func (es *EpisodesService) DeletePodcastEpisode(uid, id string) (bool, error) {
	episode, err := es.GetPodcastEpisodeById(id)
	if err != nil {
		return false, err
	}

	creatorId, err := strconv.ParseUint(uid, 0, 64)
	if err != nil || episode.CreatorId != uint(creatorId) {
		return false, errors.New("episode not found")
	}

	podcast, err := es.ps.GetPodcastById(fmt.Sprint(episode.PodcastId))
	if err != nil {
		return false, err
	}

	if podcast.CreatorId != episode.CreatorId {
		return true, errors.New("unauthorized creator, cannot delete")
	}

	if episode.Visibility != "draft" {
		return true, errors.New("cannot delete a non draft episode")
	}

	return es.er.Destroy(episode)
}
