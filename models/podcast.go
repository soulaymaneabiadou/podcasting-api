package models

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Podcast struct {
	Model
	Name        string `json:"name" gorm:"type:varchar(255);not null;unique"`
	Slug        string `json:"slug" gorm:"type:varchar(255);not null;unique"`
	Description string `json:"description" gorm:"type:varchar(255);not null;unique"`
	CreatorId   uint   `json:"creator_id" gorm:"type:varchar(255);not null;index"`
	// Picture     string      `json:"picture" gorm:"type:varchar(255)"`
	// SocialLinks StringSlice `json:"social_links" gorm:"type:varchar(255)"`
	// Hosts       StringSlice `json:"hosts" gorm:"type:varchar(255)"`
	// Tags        StringSlice `json:"tags" gorm:"type:varchar(255)"`

	Creator User `json:"creator" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (p *Podcast) BeforeSave(db *gorm.DB) error {
	p.Slug = slug.Make(p.Name)

	return nil
}

// type StringSlice []string

// func (o *StringSlice) Scan(src any) error {
// 	bytes, ok := src.([]byte)
// 	if !ok {
// 		return errors.New("src value cannot cast to []byte")
// 	}

// 	*o = strings.Split(string(bytes), ",")
// 	return nil
// }

// func (o StringSlice) Value() (driver.Value, error) {
// 	if len(o) == 0 {
// 		return nil, nil
// 	}
// 	return strings.Join(o, ","), nil
// }
