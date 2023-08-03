package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/models"
)

type UserRepoInterface interface {
	RegisterUser(user models.User) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &userRepository{db: db}
}

func (r *userRepository) RegisterUser(user models.User) (models.User, error) {
	tx := r.db.Begin()
	if err := tx.Debug().Create(&user).Error; err != nil {
		tx.Rollback()
		return user, fmt.Errorf("[Register.Insert] Error when query save data with : %w", err)
	}
	tx.Commit()
	return user, nil
}