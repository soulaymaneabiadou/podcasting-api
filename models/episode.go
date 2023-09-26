package models

import (
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Episode struct {
	Model
	Title       string      `json:"name" gorm:"type:varchar(255);not null;unique"`
	Description string      `json:"description" gorm:"type:varchar(255);not null;unique"`
	MediaLink   string      `json:"media_link" gorm:"type:varchar(255)"`
	Visibility  string      `json:"visibility" gorm:"type:visibility;default:'draft'" sql:"type:Visibility"`
	Tags        StringSlice `json:"tags" gorm:"type:varchar(255)"`
	CreatorId   uint        `json:"creator_id" gorm:"type:varchar(255);not null;index"`
	PodcastId   uint        `json:"podcast_id" gorm:"type:varchar(255);not null;index"`
	Slug        string      `json:"slug" gorm:"type:varchar(255);not null;unique"`
	PublishedAt time.Time   `json:"published_at" gorm:"null"`

	Creator User    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Podcast Podcast `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (e *Episode) BeforeSave(db *gorm.DB) error {
	e.Slug = slug.Make(e.Title)

	return nil
}
