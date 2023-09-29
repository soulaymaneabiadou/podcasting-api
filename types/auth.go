package types

type SignupInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email" gorm:"unique"`
	Password string `json:"password" binding:"required,min=8,alphanum"`
}

type SigninInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	IP       string `json:"-" binding:"-"`
}

type UpdateDetailsInput struct {
	Name  string `json:"name" binding:"omitempty"`
	Email string `json:"email" binding:"omitempty,email"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdatePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,alphanum"`
}

type ResetPasswordInput struct {
	Password string `json:"password" binding:"required,min=8,alphanum"`
}
