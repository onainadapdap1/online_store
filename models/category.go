package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	UserID      uint
	Name        string `gorm:"unique;not null"`
	Description string
	Slug        string `gorm:"unique;not null"`
	ImageURL    string
	User        User
}

func (c *Category) TableName() string {
	return "tb_categories"
}

func (c *Category) BeforeSave() (err error) {
	c.Slug = slug.Make(c.Name)
	return
}
