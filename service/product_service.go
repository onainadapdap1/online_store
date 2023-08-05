package service

import (
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/repository"
)


type ProductServiceInterface interface {
	CreateProduct(input dtos.CreateProductInput) (models.Product, error)
	GetCategoryByID(inputID uint) (models.Category, error)
	UpdateProduct(inputSlug dtos.GetProductDetailInput, inputData dtos.CreateProductInput) (models.Product, error)
	FindProductBySlug(inputSlug dtos.GetProductDetailInput) (models.Product, error)
	FindAllProduct() ([]models.Product, error)
}

type productService struct {
	repo repository.ProductRepoInterface
}

func NewProductService(repo repository.ProductRepoInterface) ProductServiceInterface {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(input dtos.CreateProductInput) (models.Product, error) {
	product := models.Product{
		UserID:      input.User.ID,
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
		ImageURL:    input.ImageURL,
	}

	product, err := s.repo.CreateProduct(product)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *productService) GetCategoryByID(inputID uint) (models.Category, error) {
	category, err := s.repo.GetCategoryByID(inputID)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (s *productService) FindProductBySlug(inputSlug dtos.GetProductDetailInput) (models.Product, error) {
	product, err := s.repo.FindProductBySlug(inputSlug.Slug)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *productService) UpdateProduct(inputSlug dtos.GetProductDetailInput, inputData dtos.CreateProductInput) (models.Product, error) {
	// product := models.Product{}
	product, err := s.repo.FindProductBySlug(inputSlug.Slug)
	if err != nil {
		return product, err
	}

	product.Name = inputData.Name
	product.Description = inputData.Description
	product.Price = inputData.Price
	product.Quantity = inputData.Quantity
	product.ImageURL = inputData.ImageURL
	product.CategoryID = inputData.CategoryID

	updatedProduct, err := s.repo.UpdateProduct(product)
	if err != nil {
		return updatedProduct, err
	}
	return updatedProduct, nil
}

func (s *productService) FindAllProduct() ([]models.Product, error) {
	products, err := s.repo.FindAllProduct()
	if err != nil {
		return products, err
	}

	return products, nil
}
