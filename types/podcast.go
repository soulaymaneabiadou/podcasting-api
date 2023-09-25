package types

import "podcast/models"

type Podcast = models.Podcast

type CreatePodcastInput struct {
	Name        string             `json:"name" binding:"required,min=2"`
	Description string             `json:"description" binding:"required,min=2"`
	Picture     string             `json:"picture" binding:"omitempty"`
	SocialLinks models.StringSlice `json:"social_links" binding:"omitempty"`
	Hosts       models.StringSlice `json:"hosts" binding:"omitempty"`
	Tags        models.StringSlice `json:"tags" binding:"omitempty"`
	CreatorId   uint               `json:"-" binding:"-"`
}

type UpdatePodcastInput struct {
	Description string             `json:"description" binding:"omitempty"`
	Picture     string             `json:"picture" binding:"omitempty"`
	SocialLinks models.StringSlice `json:"social_links" binding:"omitempty"`
	Hosts       models.StringSlice `json:"hosts" binding:"omitempty"`
	Tags        models.StringSlice `json:"tags" binding:"omitempty"`
}
