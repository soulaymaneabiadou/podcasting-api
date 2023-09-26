package types

import (
	"podcast/models"
	"time"
)

type Episode = models.Episode

type CreateEpisodeInput struct {
	Title       string             `json:"title" binding:"required,min=2"`
	Description string             `json:"description" binding:"required,min=2"`
	MediaLink   string             `json:"media_link" binding:"omitempty"`
	Visibility  string             `json:"visibility" binding:"omitempty"`
	Tags        models.StringSlice `json:"tags" binding:"omitempty"`
	PublishedAt time.Time          `json:"published_at" binding:"omitempty"`
	CreatorId   uint               `json:"-" binding:"-"`
	PodcastId   uint               `json:"-" binding:"-"`
}

type UpdateEpisodeInput struct {
	Title       string             `json:"title" binding:"omitempty"`
	Description string             `json:"description" binding:"omitempty"`
	MediaLink   string             `json:"media_link" binding:"omitempty"`
	Visibility  string             `json:"visibility" binding:"omitempty"`
	Tags        models.StringSlice `json:"tags" binding:"omitempty"`
	PublishedAt time.Time          `json:"published_at" binding:"omitempty"`
}

type PublishEpisodeInput struct {
	Visibility string `json:"visibility" binding:"required"`
}
