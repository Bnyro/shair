package entities

import "gorm.io/gorm"

type BlogPost struct {
	gorm.Model
	Title     string `json:"title"`
	Body      string `json:"body"`
	Reference string `json:"reference"`
}
