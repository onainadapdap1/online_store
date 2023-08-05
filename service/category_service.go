package service

import (
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/repository"
)

type CategoryServiceInterface interface {
	CreateCategory(input dtos.CreateCategoryInput) (models.Category, error)
	UpdateCategory(inputSlug dtos.GetCategoryDetailInput, inputData dtos.CreateCategoryInput) (models.Category, error)
	FindBySlug(inputSlug dtos.GetCategoryDetailInput) (models.Category, error)
	FindByCategoryID(categoryID uint) (models.Category, error)
	FindAllCategory() ([]models.Category, error)
	DeleteCategory(category models.Category) error
}

type categoryService struct {
	repo repository.CategoryRepoInterface
}

func NewCategoryService(repo repository.CategoryRepoInterface) CategoryServiceInterface {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(input dtos.CreateCategoryInput) (models.Category, error) {
	category := models.Category{
		UserID:      input.User.ID,
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
	}

	category, err := s.repo.CreateCategory(category)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(inputSlug dtos.GetCategoryDetailInput, inputData dtos.CreateCategoryInput) (models.Category, error) {
	category, err := s.repo.FindBySlug(inputSlug.Slug)
	if err != nil {
		return category, err
	}

	category.Name = inputData.Name
	category.Description = inputData.Description
	category.ImageURL = inputData.ImageURL

	updatedCategory, err := s.repo.UpdateCategory(category)
	if err != nil {
		return updatedCategory, err
	}

	return updatedCategory, nil
}

func (s *categoryService) FindBySlug(inputSlug dtos.GetCategoryDetailInput) (models.Category, error) {
	category, err := s.repo.FindBySlug(inputSlug.Slug)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (s *categoryService) FindByCategoryID(categoryID uint) (models.Category, error) {
	category, err := s.repo.FindByCategoryID(categoryID)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (s *categoryService) FindAllCategory() ([]models.Category, error) {
	categories, err := s.repo.FindAllCategory()
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *categoryService) DeleteCategory(category models.Category) error {
	if err := s.repo.DeleteCategory(category); err != nil {
		return err
	}
	return nil
}