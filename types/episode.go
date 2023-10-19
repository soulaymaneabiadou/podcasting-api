package types

import (
	"podcast/models"
	"time"
)

type Episode = models.Episode

type CreateEpisodeInput struct {
	Title       string             `form:"title" json:"title" binding:"required,min=2"`
	Description string             `form:"description" json:"description" binding:"required,min=2"`
	Media       string             `form:"media" json:"media" binding:"omitempty"`
	MediaLink   string             `form:"-" json:"media_link" binding:"omitempty"`
	Visibility  string             `form:"visibility" json:"visibility" binding:"omitempty"`
	Tags        models.StringSlice `form:"tags" json:"tags" binding:"omitempty"`
	PublishedAt time.Time          `form:"-,omitempty" json:"-,omitempty" binding:"omitempty"`
	CreatorId   uint               `form:"-" json:"-" binding:"-"`
	PodcastId   uint               `form:"-" json:"-" binding:"-"`
}

type UpdateEpisodeInput struct {
	Title       string             `form:"title" json:"title" binding:"omitempty"`
	Description string             `form:"description" json:"description" binding:"omitempty"`
	Media       string             `form:"media" json:"media" binding:"-"`
	MediaLink   string             `form:"-" json:"-" binding:"omitempty"`
	Visibility  string             `form:"visibility" json:"visibility" binding:"omitempty"`
	Tags        models.StringSlice `form:"tags" json:"tags" binding:"omitempty"`
	PublishedAt time.Time          `form:"-" json:"-" binding:"omitempty"`
}

type PublishEpisodeInput struct {
	Visibility string `json:"visibility" binding:"required"`
}
