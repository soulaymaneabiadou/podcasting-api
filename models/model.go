package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// equivilant to using `gorm.Model`, but with more control over the json struct
type Model struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
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

type SocialLinks map[string]string

func (o *SocialLinks) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("map `src` value cannot cast to []byte")
	}

	return json.Unmarshal(bytes, o)
}

func (o SocialLinks) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}

	return json.Marshal(o)
}
