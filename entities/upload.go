package entities

import "gorm.io/gorm"

type Upload struct {
	gorm.Model
	Name    string `json:"name"`
	Token   string `json:"token"`
	Size    int64  `json:"size"`
	Created int64  `json:"created"`
	Expires int64  `json:"expires"`
	UserID  uint   `json:"userId"`
}
