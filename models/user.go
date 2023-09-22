package models

import (
	"html"
	"strings"
	"time"

	"podcast/hasher"

	"gorm.io/gorm"
)

type User struct {
	Model
	Name                string    `json:"name" gorm:"type:varchar(255)"`
	Email               string    `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password            string    `json:"-" gorm:"type:varchar(255);not null"`
	ResetPasswordToken  string    `json:"-" gorm:"type:varchar(255);null"`
	ResetPasswordExpire time.Time `json:"-"`
}

func (u *User) BeforeSave(db *gorm.DB) error {
	var err error

	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	u.Password, err = hasher.HashPassword(u.Password)

	return err
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	var err error

	if u.Password != "" {
		u.Password, err = hasher.HashPassword(u.Password)
	}

	return err
}
