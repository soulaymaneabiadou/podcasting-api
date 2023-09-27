package models

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Podcast struct {
	Model
	Name        string      `json:"name" gorm:"type:varchar(255);not null;unique"`
	Slug        string      `json:"slug" gorm:"type:varchar(255);not null;unique"`
	Description string      `json:"description" gorm:"type:varchar(255);not null;unique"`
	Picture     string      `json:"picture" gorm:"type:varchar(255)"`
	SocialLinks SocialLinks `json:"social_links" gorm:"type:JSON"`
	Hosts       StringSlice `json:"hosts" gorm:"type:varchar(255)"`
	Tags        StringSlice `json:"tags" gorm:"type:varchar(255)"`
	CreatorId   uint        `json:"creator_id" gorm:"type:varchar(255);not null;index"`

	Creator       User           `json:"creator" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Episodes      []Episode      `json:"episodes"`
	Subscriptions []Subscription `json:"-"`
	Subscribers   []User         `json:"-" gorm:"many2many:subscriptions"`
}

func (p *Podcast) BeforeSave(db *gorm.DB) error {
	p.Slug = slug.Make(p.Name)

	return nil
}
