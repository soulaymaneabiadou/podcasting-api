package types

type SignupInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,alphanum"`
}

type SigninInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateDetailsInput struct {
	Name  string `json:"name" binding:"-"`
	Email string `json:"email" binding:"email"`
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
