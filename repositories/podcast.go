package repositories

import (
	"errors"
	"log"

	"podcast/database"
	"podcast/types"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type PodcastsRepository struct {
	db *gorm.DB
}

func NewPodcastsRepository(db *gorm.DB) *PodcastsRepository {
	return &PodcastsRepository{db: db}
}

func (ur *PodcastsRepository) Count() (int64, error) {
	var count int64

	if err := ur.db.Find(&types.Podcast{}).Count(&count).Error; err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return count, nil
}

func (ur *PodcastsRepository) GetAll(p types.Paginator) ([]types.Podcast, error) {
	var podcasts []types.Podcast

	if err := database.Paginate(p).Preload("Creator").Find(&podcasts).Error; err != nil {
		log.Println(err.Error())
		return []types.Podcast{}, err
	}

	return podcasts, nil
}

func (ur *PodcastsRepository) GetById(id string) (types.Podcast, error) {
	var podcast types.Podcast

	if err := ur.db.Preload("Creator").Preload("Episodes", "visibility IN (?)", "public").Where("id=?", id).First(&podcast).Error; err != nil {
		log.Println(err.Error())
		return types.Podcast{}, errors.New("podcast not found")
	}

	return podcast, nil
}

func (ur *PodcastsRepository) GetBySlug(slug string) (types.Podcast, error) {
	var podcast types.Podcast

	if err := ur.db.Preload("Creator").Preload("Episodes", "visibility IN (?)", "public").Where("slug=?", slug).First(&podcast).Error; err != nil {
		log.Println(err, podcast)
		return types.Podcast{}, err
	}

	return podcast, nil
}

func (ur *PodcastsRepository) Create(input types.CreatePodcastInput) (types.Podcast, error) {
	podcast := types.Podcast{
		Name:        input.Name,
		Headline:    input.Headline,
		Description: input.Description,
		CreatorId:   input.CreatorId,
		Picture:     input.Picture,
		SocialLinks: input.SocialLinks,
		Hosts:       input.Hosts,
		Tags:        input.Tags,
	}

	if err := ur.db.Create(&podcast).Error; err != nil {
		log.Println(err.Error())
		return types.Podcast{}, err
	}

	return podcast, nil
}

func (ur *PodcastsRepository) Update(podcast types.Podcast, input types.UpdatePodcastInput) (types.Podcast, error) {
	payload := types.Podcast{
		Description: input.Description,
		Picture:     input.Picture,
		SocialLinks: input.SocialLinks,
		Hosts:       input.Hosts,
		Tags:        input.Tags,
		Headline:    input.Headline,
	}

	if payload.Name != "" {
		payload.Slug = slug.Make(payload.Name)
	}

	if err := ur.db.Model(&podcast).Updates(payload).Error; err != nil {
		log.Println(err.Error())
		return podcast, err
	}

	return podcast, nil
}

func (ur *PodcastsRepository) Destroy(podcast types.Podcast) (bool, error) {
	if err := ur.db.Delete(&podcast).Error; err != nil {
		return false, err
	}

	return true, nil
}
