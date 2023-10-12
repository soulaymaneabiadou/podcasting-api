package repositories

import (
	"errors"
	"log"
	"strings"

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

func (ur *PodcastsRepository) GetAll(f types.PodcastFilters, s types.Sorter, p types.Paginator) ([]types.Podcast, error) {
	var podcasts []types.Podcast

	db := filterPodcasts(ur.db, f)
	db = database.Sort(db, s)
	if err := database.Paginate(db, p).Preload("Creator").Find(&podcasts).Error; err != nil {
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

func (ur *PodcastsRepository) GetByCreatorId(id string) (types.Podcast, error) {
	var podcast types.Podcast

	if err := ur.db.Where("creator_id=?", id).First(&podcast).Error; err != nil {
		log.Println(err.Error())
		return types.Podcast{}, errors.New("podcast not found")
	}

	return podcast, nil
}

func (pr *PodcastsRepository) GetSubscriptions(pid string) ([]types.Subscription, error) {
	var subscriptions []types.Subscription

	if err := pr.db.Where("podcast_id=?", pid).Find(&subscriptions).Error; err != nil {
		return []types.Subscription{}, err
	}

	return subscriptions, nil
}

func (pr *PodcastsRepository) GetSubscriptionsCount(pid string) (int64, error) {
	var count int64

	if err := pr.db.Model(&types.Subscription{}).Where("podcast_id=?", pid).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *PodcastsRepository) GetActiveSubscriptionsCount(pid string) (int64, error) {
	var count int64

	if err := pr.db.Model(&types.Subscription{}).Where("status=?", "active").Where("podcast_id=?", pid).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *PodcastsRepository) GetEpisodes(pid string) ([]types.Episode, error) {
	var episodes []types.Episode

	if err := pr.db.Where("podcast_id=?", pid).Find(&episodes).Error; err != nil {
		return []types.Episode{}, err
	}

	return episodes, nil
}

func (pr *PodcastsRepository) GetEpisodesCount(pid string) (int64, error) {
	var count int64

	if err := pr.db.Model(&types.Episode{}).Where("podcast_id=?", pid).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *PodcastsRepository) GetByListenerId(id string) ([]types.Podcast, error) {
	var subscriptions []types.Subscription
	var podcasts []types.Podcast = []types.Podcast{}

	if err := pr.db.Where("user_id=?", id).Find(&subscriptions).Error; err != nil {
		return podcasts, err
	}

	for _, sub := range subscriptions {
		podcast := types.Podcast{}
		if err := pr.db.First(&podcast, sub.PodcastId).Error; err != nil {
			return podcasts, err
		}

		podcasts = append(podcasts, podcast)
	}

	return podcasts, nil
}

func filterPodcasts(db *gorm.DB, f types.PodcastFilters) *gorm.DB {
	return db.Scopes(func(db *gorm.DB) *gorm.DB {
		query := db

		if f.Name != "" {
			query = query.Where("lower(name) LIKE ?", "%"+strings.ToLower(f.Name)+"%")
		}
		if f.Host != "" {
			query = query.Where("lower(hosts) LIKE ?", "%"+strings.ToLower(f.Host)+"%")
		}

		return query
	})
}
