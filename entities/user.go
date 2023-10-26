package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string `json:"username"`
	Password   string `json:"-"`
	TotpSecret string `json:"-"`
	AuthToken  string `json:"auth"`
}
