package models

import (
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
