package entities

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Token       string `json:"token"`
}
