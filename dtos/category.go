package dtos

import "github.com/onainadapdap1/online_store/models"

/* CATEGORY INPUT */

// create category input
type CreateCategoryInput struct {
	Name        string `gorm:"not null" form:"name" json:"name"`
	Description string `gorm:"not null" form:"description" json:"description"`
	ImageURL    string `gorm:"not null" form:"image" json:"image"`
	User        models.User
}

type GetCategoryDetailInput struct {
	Slug string `uri:"slug" binding:"required"`
}

/* END CATEGORY INPUT */

/*CATEGORY*/
// format ketika insert category
type CategoryFormatter struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
	Slug        string `json:"slug" gorm:"uniqueIndex;not null"`
	ImageURL    string `json:"image_url"`
}

func FormatCategory(category models.Category) CategoryFormatter {
	categoryFormatter := CategoryFormatter{
		ID:          category.ID,
		UserID:      category.UserID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		ImageURL:    category.ImageURL,
	}

	return categoryFormatter
}

type CategoryDetailFormatter struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	Name        string                `json:"category_name"`
	Description string                `json:"description"`
	Slug        string                `json:"slug"`
	ImageURL    string                `json:"image_url"`
	User        CategoryUserFormatter `json:"user"`
}

type CategoryUserFormatter struct {
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func FormateCategoryDetail(category models.Category) CategoryDetailFormatter {
	categoryDetailFormatter := CategoryDetailFormatter{
		ID:          category.ID,
		UserID:      category.UserID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		ImageURL:    category.ImageURL,
	}
	user := category.User

	categoryUserFormatter := CategoryUserFormatter{
		FullName: user.FullName,
		Role:     user.Role,
	}

	categoryDetailFormatter.User = categoryUserFormatter

	return categoryDetailFormatter
}

func FormatCategories(categories []models.Category) []CategoryDetailFormatter {
	categoriesFormatter := []CategoryDetailFormatter{}

	for _, category := range categories {
		categoryFormatter := FormateCategoryDetail(category)
		categoriesFormatter = append(categoriesFormatter, categoryFormatter)
	}

	return categoriesFormatter
}
