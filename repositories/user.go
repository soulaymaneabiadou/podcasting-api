package repositories

import (
	"errors"
	"log"
	"time"

	"podcast/database"
	"podcast/hasher"
	"podcast/types"

	"gorm.io/gorm"
)

type UsersRepository struct {
	database *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{database: db}
}

func (ur *UsersRepository) GetAll(p database.Paginator) ([]types.User, error) {
	var users []types.User

	if err := database.Paginate(p).Find(&users).Error; err != nil {
		log.Println(err)
		return []types.User{}, err
	}

	return users, nil
}

func (ur *UsersRepository) GetById(id string) (types.User, error) {
	var user types.User

	if err := database.DB.Where("id=?", id).First(&user).Error; err != nil {
		log.Println(err)
		return types.User{}, errors.New("user not found")
	}

	return user, nil
}

func (ur *UsersRepository) GetByEmail(email string) (types.User, error) {
	var user types.User

	if err := database.DB.Where("email=?", email).First(&user).Error; err != nil {
		log.Println(err, user)
		return types.User{}, err
	}

	return user, nil
}

func (ur *UsersRepository) GetByResetPasswordToken(token string) (types.User, error) {
	var user types.User

	if err := database.DB.Where("reset_password_token=?", token).Where("reset_password_expire>?", time.Now()).First(&user).Error; err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (ur *UsersRepository) Create(u types.CreateUserInput) (types.User, error) {
	user := types.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     string(u.Role),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Println(err)
		return types.User{}, err
	}

	return user, nil
}

func (ur *UsersRepository) Update(user types.User, input types.UpdateUserInput) (types.User, error) {
	var err error

	payload := types.User{
		Name:                input.Name,
		Email:               input.Email,
		Password:            input.Password,
		ResetPasswordToken:  input.ResetPasswordToken,
		ResetPasswordExpire: input.ResetPasswordExpire,
	}

	if payload.Password != "" {
		payload.Password, err = hasher.HashPassword(payload.Password)
		if err != nil {
			return user, err
		}
	}

	if err := database.DB.Model(&user).Updates(payload).Error; err != nil {
		log.Println(err)
		return user, err
	}

	return user, nil
}

func (ur *UsersRepository) Destroy(user types.User) (bool, error) {
	if err := database.DB.Delete(&user).Error; err != nil {
		return false, err
	}

	return true, nil
}
