package models

import (
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Episode struct {
	Model
	Title       string      `json:"title" gorm:"type:varchar(255);not null;unique"`
	Slug        string      `json:"slug" gorm:"type:varchar(255);not null;unique"`
	Description string      `json:"description" gorm:"type:varchar(255);not null;unique"`
	MediaLink   string      `json:"media_link" gorm:"type:varchar(255)"`
	Visibility  string      `json:"visibility" gorm:"type:visibility;default:'draft'" sql:"type:Visibility"`
	Tags        StringSlice `json:"tags,omitempty" gorm:"type:varchar(255)"`
	PublishedAt time.Time   `json:"published_at,omitempty" gorm:"default:null"`
	CreatorId   uint        `json:"-" gorm:"not null;index"`
	PodcastId   uint        `json:"-" gorm:"not null;index"`

	Creator User    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Podcast Podcast `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (e *Episode) BeforeSave(db *gorm.DB) error {
	e.Slug = slug.Make(e.Title)

	return nil
}
