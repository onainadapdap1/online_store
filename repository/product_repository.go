package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/models"
)

type ProductRepoInterface interface {
	CreateProduct(product models.Product) (models.Product, error)
	UpdateProduct(product models.Product) (models.Product, error)
	GetCategoryByID(id uint) (models.Category, error)
	FindProductBySlug(slug string) (models.Product, error)
	FindAllProduct() ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepoInterface {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product models.Product) (models.Product, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Create(&product).Error; err != nil {
		tx.Rollback()
		return product, fmt.Errorf("[CreateProduct.Insert] Error when query save data with : %w", err)
	}

	tx.Commit()

	return product, nil
}

func (r *productRepository) GetCategoryByID(id uint) (models.Category, error) {
	var category models.Category
	tx := r.db.Begin()

	// SELECT c.*, u.*
	// FROM categories c
	// JOIN users u ON u.id = c.user_id
	// WHERE c.id = ?;
	err := tx.Debug().Preload("User").Where("id = ?", id).Find(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *productRepository) FindProductBySlug(slug string) (models.Product, error) {
	tx := r.db.Begin()
	var product models.Product
	err := tx.Debug().Preload("User").Preload("Category").Where("slug = ?", slug).Find(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *productRepository) UpdateProduct(product models.Product) (models.Product, error) {
	tx := r.db.Begin()

	err := tx.Debug().Save(&product).Error
	if err != nil {
		tx.Rollback()
		return product, err
	}
	tx.Commit()
	return product, nil
}

func (r *productRepository) FindAllProduct() ([]models.Product, error) {
	tx := r.db.Begin()
	products := []models.Product{}

	// SELECT *
	// FROM products
	// LEFT JOIN users ON products.user_id = users.id
	// LEFT JOIN categories ON products.category_id = categories.id;
	err := tx.Debug().Preload("User").Preload("Category").Find(&products).Error
	if err != nil {
		return products, err
	}

	return products, nil
}
