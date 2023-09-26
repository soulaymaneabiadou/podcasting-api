package repositories

import (
	"errors"
	"log"

	"podcast/database"
	"podcast/types"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type EpisodesRepository struct {
	db *gorm.DB
}

func NewEpisodesRepository(db *gorm.DB) *EpisodesRepository {
	return &EpisodesRepository{db: db}
}

func (er *EpisodesRepository) Count(pid string) (int64, error) {
	var count int64

	if err := er.db.Where("podcast_id=?", pid).Find(&types.Episode{}).Count(&count).Error; err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return count, nil
}

func (er *EpisodesRepository) GetAll(pid string, p types.Paginator) ([]types.Episode, error) {
	var episodes []types.Episode

	if err := database.Paginate(p).Where("podcast_id=?", pid).Where("visibility=?", "public").Find(&episodes).Error; err != nil {
		log.Println(err.Error())
		return []types.Episode{}, err
	}

	return episodes, nil
}

func (er *EpisodesRepository) GetById(id string) (types.Episode, error) {
	var episode types.Episode

	if err := er.db.Where("id=?", id).First(&episode).Error; err != nil {
		log.Println(err.Error())
		return types.Episode{}, errors.New("podcast episode not found")
	}

	return episode, nil
}

func (er *EpisodesRepository) GetBySlug(slug string) (types.Episode, error) {
	var episode types.Episode

	if err := er.db.Where("slug=?", slug).First(&episode).Error; err != nil {
		log.Println(err, episode)
		return types.Episode{}, err
	}

	return episode, nil
}

func (er *EpisodesRepository) Create(input types.CreateEpisodeInput) (types.Episode, error) {
	episode := types.Episode{
		Title:       input.Title,
		Description: input.Description,
		MediaLink:   input.MediaLink,
		Tags:        input.Tags,
		PodcastId:   input.PodcastId,
		CreatorId:   input.CreatorId,
		Visibility:  input.Visibility,
		PublishedAt: input.PublishedAt,
	}

	if err := er.db.Create(&episode).Error; err != nil {
		log.Println(err.Error())
		return types.Episode{}, err
	}

	return episode, nil
}

func (er *EpisodesRepository) Update(episode types.Episode, input types.UpdateEpisodeInput) (types.Episode, error) {
	payload := types.Episode{
		Title:       input.Title,
		Description: input.Description,
		MediaLink:   input.MediaLink,
		Tags:        input.Tags,
		Visibility:  input.Visibility,
		PublishedAt: input.PublishedAt,
	}

	if payload.Title != "" {
		payload.Slug = slug.Make(payload.Title)
	}

	if err := er.db.Model(&episode).Updates(payload).Error; err != nil {
		log.Println(err.Error())
		return episode, err
	}

	return episode, nil
}

func (er *EpisodesRepository) Destroy(episode types.Episode) (bool, error) {
	if err := er.db.Delete(&episode).Error; err != nil {
		return false, err
	}

	return true, nil
}
