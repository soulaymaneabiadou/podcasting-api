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
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (ur *UsersRepository) GetAll(p types.Paginator) ([]types.User, error) {
	var users []types.User

	if err := database.Paginate(ur.db, p).Find(&users).Error; err != nil {
		log.Println(err)
		return []types.User{}, err
	}

	return users, nil
}

func (ur *UsersRepository) GetById(id string) (types.User, error) {
	var user types.User

	if err := ur.db.Where("id=?", id).Where("verified=?", true).First(&user).Error; err != nil {
		log.Println(err)
		return types.User{}, errors.New("user not found")
	}

	return user, nil
}

func (ur *UsersRepository) GetByStripeAccountId(aid string) (types.User, error) {
	var user types.User

	if err := ur.db.Where("stripe_account_id=?", aid).First(&user).Error; err != nil {
		log.Println(err)
		return types.User{}, errors.New("user not found by account id")
	}

	return user, nil
}

func (ur *UsersRepository) GetByEmail(email string) (types.User, error) {
	var user types.User

	if err := ur.db.Where("email=?", email).First(&user).Error; err != nil {
		log.Println(err, user)
		return types.User{}, err
	}

	return user, nil
}

func (ur *UsersRepository) GetByResetPasswordToken(token string) (types.User, error) {
	var user types.User

	if err := ur.db.Where("reset_password_token=?", token).Where("reset_password_expire>?", time.Now()).First(&user).Error; err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (ur *UsersRepository) GetByVerificationToken(token string) (types.User, error) {
	var user types.User

	if err := ur.db.Where("verification_token=?", token).Where("verified=?", false).First(&user).Error; err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (ur *UsersRepository) Create(u types.CreateUserInput) (types.User, error) {
	user := types.User{
		Name:              u.Name,
		Email:             u.Email,
		Password:          u.Password,
		Role:              string(u.Role),
		VerificationToken: u.VerificationToken,
	}

	if err := ur.db.Create(&user).Error; err != nil {
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
		StripeCustomerId:    input.StripeCustomerId,
		StripeAccountId:     input.StripeAccountId,
		ChargesEnabled:      input.ChargesEnabled,
		PayoutsEnabled:      input.PayoutsEnabled,
		DetailsSubmitted:    input.DetailsSubmitted,
		VerificationToken:   input.VerificationToken,
		Verified:            input.Verified,
		VerifiedAt:          input.VerifiedAt,
		SigninCount:         input.SigninCount,
		CurrentSigninAt:     input.CurrentSigninAt,
		CurrentSigninIP:     input.CurrentSigninIP,
		LastSigninAt:        input.LastSigninAt,
		LastSigninIP:        input.LastSigninIP,
	}

	if payload.Password != "" {
		payload.Password, err = hasher.HashPassword(payload.Password)
		if err != nil {
			return user, err
		}
	}

	if err := ur.db.Model(&user).Updates(payload).Error; err != nil {
		log.Println(err)
		return user, err
	}

	return user, nil
}

func (ur *UsersRepository) Destroy(user types.User) (bool, error) {
	if err := ur.db.Delete(&user).Error; err != nil {
		return false, err
	}

	return true, nil
}
