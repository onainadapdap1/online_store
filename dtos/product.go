package dtos

import (
	"github.com/gosimple/slug"
	"github.com/onainadapdap1/online_store/models"
)

/* PRODUCT INPUT */
type CreateProductInput struct {
	Name        string  `gorm:"not null" form:"name" json:"name"`
	Description string  `gorm:"not null" form:"description" json:"description"`
	Price       float64 `gorm:"not null" form:"price" json:"price"`
	Quantity    int     `gorm:"not null" form:"quantity" json:"quantity"`
	ImageURL    string  `gorm:"not null" form:"image" json:"image"`
	CategoryID  uint    `gorm:"not null" form:"category_id" json:"category_id"`
	User        models.User
	Category    models.Category
}

type GetProductDetailInput struct {
	Slug string `uri:"slug" binding:"required"`
}

/* END PRODUCT INPUT */

/* PRODUCT */
type ProductFormatter struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	Name        string  `json:"product_name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

func FormatProduct(product models.Product) ProductFormatter {
	productFormatter := ProductFormatter{
		ID:          product.ID,
		UserID:      product.UserID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Slug:        slug.Make(product.Name),
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		ImageURL:    product.ImageURL,
	}

	return productFormatter
}

type ProductDetailFormatter struct {
	ID          uint                     `json:"id"`
	Name        string                   `json:"product_name"`
	Slug        string                   `json:"slug"`
	Description string                   `json:"description"`
	Price       float64                  `json:"price"`
	Quantity    int                      `json:"quantity"`
	ImageURL    string                   `json:"image_url"`
	UserID      uint                     `json:"user_id"`
	User        ProductUserFormatter     `json:"user"`
	CategoryID  uint                     `json:"category_id"`
	Category    ProductCategoryFormatter `json:"category"`
}

type ProductCategoryFormatter struct {
	ID          uint   `json:"category_id"`
	Name        string `json:"product_name"`
	Description string `json:"description"`
}
type ProductUserFormatter struct {
	ID       uint   `json:"user_id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func FormatProductDetail(product models.Product) ProductDetailFormatter {
	productDetailFormatter := ProductDetailFormatter{
		ID:          product.ID,
		UserID:      product.UserID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CategoryID:  product.CategoryID,
		ImageURL:    product.ImageURL,
	}
	user := product.User
	productUserFormatter := ProductUserFormatter{
		ID:       user.ID,
		FullName: user.FullName,
		Role:     user.Role,
	}
	productDetailFormatter.User = productUserFormatter

	category := product.Category
	productCategoryFormatter := ProductCategoryFormatter{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	productDetailFormatter.Category = productCategoryFormatter

	return productDetailFormatter
}

func FormatProducts(products []models.Product) []ProductDetailFormatter {
	productsFormatter := []ProductDetailFormatter{}

	for _, product := range products {
		productFormatter := FormatProductDetail(product)
		productsFormatter = append(productsFormatter, productFormatter)
	}

	return productsFormatter
}
