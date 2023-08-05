package service

import (
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/repository"
)

type PaymentCategoryServiceInterface interface {
	CreatePaymentCategory(input dtos.CreatePaymentCategoryInput) (models.PaymentCategory, error)
	GetPaymentCategoryBySlug(slug string) (models.PaymentCategory, error)
	GetPaymentCategoryByID(id uint) (models.PaymentCategory, error)
}

type paymentCategoryService struct {
	repo repository.PaymentCategoryRepositoryInterface
}

func NewPaymentCategoryService(repo repository.PaymentCategoryRepositoryInterface) PaymentCategoryServiceInterface {
	return &paymentCategoryService{repo: repo}
}

func (s *paymentCategoryService) CreatePaymentCategory(input dtos.CreatePaymentCategoryInput) (models.PaymentCategory, error) {
	paymentCategory := models.PaymentCategory {
		CategoryName: input.CategoryName,
		UserID: input.User.ID,
	}

	paymentCategory, err := s.repo.CreatePaymentCategory(paymentCategory)
	if err != nil {
		return paymentCategory, err
	}
	return paymentCategory, nil
}

func (s *paymentCategoryService) GetPaymentCategoryBySlug(slug string) (models.PaymentCategory, error) {
	paymentCategory, err := s.repo.GetPaymentCategoryBySlug(slug)
	if err != nil {
		return paymentCategory, err
	}

	return paymentCategory, nil
}

func (s *paymentCategoryService) GetPaymentCategoryByID(id uint) (models.PaymentCategory, error) {
	paymentCategory, err := s.repo.GetPaymentCategoryByID(id)
	if err != nil {
		return paymentCategory, err
	}
	return paymentCategory, nil
}