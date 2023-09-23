package types

import (
	"podcast/models"
	"time"
)

type User = models.User

type Role string

const (
	CREATOR_ROLE  Role = "creator"
	LISTENER_ROLE Role = "listener"
)

type CreateUserInput struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,alphanum"`
	Role     Role   `json:"role"`
}

type UpdateUserInput struct {
	Name                string    `json:"name" binding:"min=2"`
	Email               string    `json:"email" binding:"email"`
	Password            string    `json:"password" binding:"min=8,alphanum"`
	ResetPasswordToken  string    `json:"-"`
	ResetPasswordExpire time.Time `json:"-"`
}
