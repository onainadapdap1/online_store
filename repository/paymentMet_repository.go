package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/models"
)

type PaymentMethodRepositoryInterface interface {
	CreatePaymentMethod(paymentMethod models.PaymentMethod) (models.PaymentMethod, error)
}

type paymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepositoryInterface {
	return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) CreatePaymentMethod(paymentMethod models.PaymentMethod) (models.PaymentMethod, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Create(&paymentMethod).Error; err != nil {
		tx.Rollback()
		return paymentMethod, err
	}
	if err := tx.Debug().Where("id = ?", paymentMethod.ID).Preload("User").Preload("PaymentCategory").First(&paymentMethod).Error; err != nil {
		tx.Rollback()
		return paymentMethod, err
	}
	tx.Commit()

	return paymentMethod, nil
}