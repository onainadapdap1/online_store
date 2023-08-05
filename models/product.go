package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
	Description string
	Price float64
	Quantity int
	Slug string `gorm:"unique;not null"`
	ImageURL string
	CategoryID uint
	Category Category
	UserID uint
	User User
}

func (p *Product) TableName() string {
	return "tb_products"
}

func (p *Product) BeforeSave() (err error) {
	p.Slug = slug.Make(p.Name)
	return
}
