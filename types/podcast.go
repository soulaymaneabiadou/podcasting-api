package types

import "podcast/models"

type Podcast = models.Podcast

type SocialLinks = models.SocialLinks

type StringSlice = models.StringSlice

type CreatePodcastInput struct {
	Name        string      `form:"name" json:"name" binding:"required,min=2"`
	Headline    string      `form:"headline" json:"headline" binding:"required,min=2"`
	Description string      `form:"description" json:"description" binding:"required,min=2"`
	Picture     string      `form:"picture" json:"picture" binding:"required"`
	SocialLinks SocialLinks `form:"social_links" json:"social_links" binding:"omitempty"`
	Hosts       StringSlice `form:"hosts" json:"hosts" binding:"omitempty"`
	Tags        StringSlice `form:"tags" json:"tags" binding:"omitempty"`
	CreatorId   uint        `form:"-" json:"-" binding:"-"`
}

type UpdatePodcastInput struct {
	Headline    string      `form:"headline" json:"headline" binding:"omitempty"`
	Description string      `form:"description" json:"description" binding:"omitempty"`
	Picture     string      `form:"picture" json:"picture" binding:"omitempty"`
	SocialLinks SocialLinks `form:"social_links" json:"social_links" binding:"omitempty"`
	Hosts       StringSlice `form:"hosts" json:"hosts" binding:"omitempty"`
	Tags        StringSlice `form:"tags" json:"tags" binding:"omitempty"`
}

type PodcastStats struct {
	TotalSubsCount  int64 `json:"total_subscriptions_count,omitempty"`
	ActiveSubsCount int64 `json:"active_subscriptions_count,omitempty"`
	EpisodesCount   int64 `json:"episodes_count,omitempty"`
}

type PodcastFilters struct {
	Name string
	Host string
}
