package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/models"
)

type UserRepoInterface interface {
	RegisterUser(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
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

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	tx := r.db.Begin()
	var user models.User
	if err := tx.Debug().Where("email = ?", email).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}