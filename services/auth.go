package services

import (
	"errors"
	"log"
	"time"

	"podcast/hasher"
	"podcast/repositories"
	"podcast/types"
)

type AuthService struct {
	ur *repositories.UsersRepository
	es *EmailService
}

func NewAuthService(ur *repositories.UsersRepository, es *EmailService) *AuthService {
	return &AuthService{ur: ur, es: es}
}

func (as *AuthService) Signup(u types.SignupInput) (types.User, error) {
	user, err := as.ur.GetByEmail(u.Email)
	if user.ID != 0 || err == nil {
		return types.User{}, errors.New("email already exists")
	}

	data := types.CreateUserInput{Name: u.Name, Email: u.Email, Password: u.Password}

	user, err = as.ur.Create(data)
	if err != nil {
		return types.User{}, errors.New("an error occured, please try again later with valid information")
	}

	go as.es.SendWelcomeEmail(user)

	return user, nil
}

func (as *AuthService) Signin(u types.SigninInput) (types.User, error) {
	user, err := as.ur.GetByEmail(u.Email)
	if err != nil {
		log.Println(err)
		return types.User{}, errors.New("user not found")
	}

	err = hasher.VerifyPassword(u.Password, user.Password)
	if err != nil {
		log.Println(err)
		return types.User{}, err
	}

	return user, nil
}

func (as *AuthService) GetUser(id string) (types.User, error) {
	user, err := as.ur.GetById(id)
	if err != nil {
		return types.User{}, errors.New("the user does not exist")
	}

	return user, nil
}

func (as *AuthService) UpdateDetails(id string, input types.UpdateDetailsInput) (types.User, error) {
	user, err := as.ur.GetById(id)
	if err != nil {
		return user, err
	}

	user, err = as.ur.Update(user, types.UpdateUserInput{Email: input.Email})
	if err != nil {
		log.Println(err)
		return types.User{}, err
	}

	return user, nil
}

func (as *AuthService) UpdatePassword(id string, input types.UpdatePasswordInput) (bool, error) {
	user, _ := as.ur.GetById(id)

	err := hasher.VerifyPassword(input.CurrentPassword, user.Password)
	if err != nil {
		log.Println(err)
		return false, err
	}

	user, err = as.ur.Update(user, types.UpdateUserInput{Password: input.NewPassword})
	if err != nil {
		log.Println(err)
		return false, err
	}

	go as.es.SendPasswordUpdatedEmail(user)

	return true, nil
}

func (as *AuthService) ForgotPassword(input types.ForgotPasswordInput) (string, error) {
	user, err := as.ur.GetByEmail(input.Email)
	if err != nil {
		log.Println(err)
		return "", err
	}

	token, err := hasher.GenerateSecureToken(20)
	if err != nil {
		log.Println(err)
		return "", err
	}

	hash := hasher.GenerateTokenHash(token)

	as.ur.Update(user, types.UpdateUserInput{
		ResetPasswordToken:  hash,
		ResetPasswordExpire: time.Now().Add(time.Minute * 10),
	})

	go as.es.SendPasswordResetTokenEmail(user, token)

	return token, nil
}

func (as *AuthService) ResetPassword(token string, password string) (bool, error) {
	hash := hasher.GenerateTokenHash(token)

	user, err := as.ur.GetByResetPasswordToken(hash)
	if err != nil {
		log.Println(err)
		return false, err
	}

	as.ur.Update(user, types.UpdateUserInput{
		Password:            password,
		ResetPasswordToken:  "nil",
		ResetPasswordExpire: time.Now(),
	})

	go as.es.SendPasswordResettedEmail(user)

	return true, nil
}
