package entities

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID uint   `json:"userId"`
}
