package db

import (
	"errors"

	"github.com/shair/entities"
	"github.com/shair/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Init() {
	var err error
	Database, err = gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	if err != nil {
		panic("failed to create Database")
	}

	Database.AutoMigrate(&entities.User{}, &entities.Note{}, &entities.Upload{}, &entities.Paste{})
}

func CreateUser(username string, password string) (entities.User, error) {
	user := entities.User{
		Username: username,
	}
	if Database.Where(&user).Find(&user).RowsAffected != 0 {
		return user, errors.New("Username already taken!")
	}

	user.Password = util.HashPassword(password)
	user.AuthToken = util.GenerateSecureToken(30)

	Database.Create(&user)
	return user, nil
}

func LoginUser(username string, password string) (entities.User, error) {
	user := entities.User{
		Username: username,
	}

	if Database.Where(&user).Find(&user).RowsAffected == 0 {
		return user, errors.New("User not found!")
	}

	if !util.CheckPasswordHash(password, user.Password) {
		return user, errors.New("Invalid password!")
	}

	return user, nil
}

func FindUser(token string) (entities.User, error) {
	user := entities.User{
		AuthToken: token,
	}

	if Database.Where(&user).Find(&user).RowsAffected == 0 {
		return user, errors.New("User not found!")
	}

	return user, nil
}

func DeleteUser(user entities.User) {
	Database.Delete(&user)
}
