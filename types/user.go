package types

import (
	"podcast/models"
	"time"
)

type User = models.User

type CreateUserInput struct {
	Name     string `json:"name" binding:"required,name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,alphanum"`
}

type UpdateUserInput struct {
	Name                string    `json:"name" binding:"name"`
	Email               string    `json:"email" binding:"email"`
	Password            string    `json:"password" binding:"min=8,alphanum"`
	ResetPasswordToken  string    `json:"-"`
	ResetPasswordExpire time.Time `json:"-"`
}
