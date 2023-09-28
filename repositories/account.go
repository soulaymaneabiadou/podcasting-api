package repositories

import (
	"errors"
	"log"
	"podcast/types"

	"gorm.io/gorm"
)

type AccountsRepository struct {
	db *gorm.DB
}

func NewAccountsRepository(db *gorm.DB) *AccountsRepository {
	return &AccountsRepository{db: db}
}

func (ar *AccountsRepository) GetByUserId(uid string) (types.Account, error) {
	var account types.Account

	if err := ar.db.Where("user_id=?", uid).First(&account).Error; err != nil {
		log.Println(err)
		return types.Account{}, errors.New("user connect account not found")
	}

	return account, nil
}
