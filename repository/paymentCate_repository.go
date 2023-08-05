package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/models"
)

type PaymentCategoryRepositoryInterface interface {
	CreatePaymentCategory(paymentCategory models.PaymentCategory) (models.PaymentCategory, error)
	GetPaymentCategoryBySlug(slug string) (models.PaymentCategory, error)
	GetPaymentCategoryByID(id uint) (models.PaymentCategory, error)
}

type paymentCategoryRepository struct {
	db *gorm.DB
}

func NewPaymentCategoryRepository(db *gorm.DB) PaymentCategoryRepositoryInterface {
	return &paymentCategoryRepository{db: db}
}

func (r *paymentCategoryRepository) CreatePaymentCategory(paymentCategory models.PaymentCategory) (models.PaymentCategory, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Create(&paymentCategory).Error; err != nil {
		tx.Rollback()
		return paymentCategory, err
		// return paymentCategory, fmt.Errorf("payment category with name %s already exists", paymentCategory.CategoryName)
	}

	if err := tx.Debug().Where("slug = ?", paymentCategory.Slug).Preload("User").First(&paymentCategory).Error; err != nil {
		tx.Rollback()
		return paymentCategory, err
	}
	tx.Commit()

	return paymentCategory, nil
}

func (r *paymentCategoryRepository) GetPaymentCategoryBySlug(slug string) (models.PaymentCategory, error) {
	tx := r.db.Begin()
	var paymentCategory models.PaymentCategory
	if err := tx.Debug().Where("slug = ?", slug).Preload("User").First(&paymentCategory).Error; err != nil {
		tx.Rollback()
		return paymentCategory, err
	}
	tx.Commit()

	return paymentCategory, nil
}

func (r *paymentCategoryRepository) GetPaymentCategoryByID(id uint) (models.PaymentCategory, error) {
	tx := r.db.Begin()
	var paymentCategory models.PaymentCategory
	if err := tx.Debug().Where("id = ?", id).Preload("User").First(&paymentCategory).Error; err != nil {
		tx.Rollback()
		return paymentCategory, err
	}
	tx.Commit()
	return paymentCategory, nil
}