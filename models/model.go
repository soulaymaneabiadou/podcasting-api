package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

// equivilant to using `gorm.Model`, but with more control over the json struct
type Model struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type StringSlice []string

func (o *StringSlice) Scan(src any) error {
	bytes := []byte(src.(string))
	*o = strings.Split(string(bytes), ", ")

	return nil
}

func (o StringSlice) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}

	return strings.Join(o, ", "), nil
}

type SocialLinks struct {
	Instagram string `json:"instagram"`
	Twitter   string `json:"twitter"`
}

func (o *SocialLinks) Scan(value interface{}) error {
	bytes := value.([]byte)
	result := SocialLinks{}
	err := json.Unmarshal(bytes, &result)
	*o = SocialLinks(result)

	return err
}

func (o SocialLinks) Value() (driver.Value, error) {
	return json.Marshal(o)
}
