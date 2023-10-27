package entities

import "gorm.io/gorm"

type Paste struct {
	gorm.Model
	Title   string `json:"title"`
	Body    string `json:"body"`
	Token   string `json:"token"`
	Created int64  `json:"created"`
	Expires int64  `json:"expires"`
}
