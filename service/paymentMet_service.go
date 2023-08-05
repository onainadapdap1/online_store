package service

import (
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/repository"
)


type PaymentMethodServiceInterface interface {
	CreatePaymentMethod(input dtos.CreatePaymentMethodInput) (models.PaymentMethod, error)
}

type paymentMethodService struct {
	repo repository.PaymentMethodRepositoryInterface
}

func NewPaymentMethodService(repo repository.PaymentMethodRepositoryInterface) PaymentMethodServiceInterface {
	return &paymentMethodService{repo: repo}
}

func (s *paymentMethodService) CreatePaymentMethod(input dtos.CreatePaymentMethodInput) (models.PaymentMethod, error) {
	paymentMethod := models.PaymentMethod {
		UserID: input.User.ID,
		CategoryPaymentID: input.CategoryPaymentID,
		MethodName: input.MethodName,
		Number: input.Number,
		OwnerName: input.OwnerName,
		CategoryName: input.PaymentCategory.CategoryName,
	}
	paymentMethod, err := s.repo.CreatePaymentMethod(paymentMethod)
	if err != nil {
		return paymentMethod, err
	}

	return paymentMethod, nil
}