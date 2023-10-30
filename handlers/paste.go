package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shair/config"
	"github.com/shair/db"
	"github.com/shair/entities"
	"github.com/shair/util"
)

func NewPaste(c echo.Context) error {
	if c.Request().Method == "GET" {
		return c.Render(http.StatusOK, "newpaste.html", nil)
	}

	now := time.Now().Unix()
	paste := entities.Paste{
		Title:   c.FormValue("title"),
		Body:    c.FormValue("body"),
		Created: now,
		Expires: now + config.TempLifeTime,
		Token:   util.GenerateSecureToken(30),
	}

	db.Database.Create(&paste)

	return util.CreateSuccessResult(c, "paste", paste.Token)
}

func GetPaste(c echo.Context) error {
	paste := entities.Paste{
		Token: c.Param("token"),
	}

	if db.Database.Where(&paste).Find(&paste).RowsAffected == 0 {
		c.Render(http.StatusNotFound, "notfound.html", nil)
	}

	return c.Render(http.StatusOK, "paste.html", paste)
}

func DeleteExpiredPastes() {
	currentTime := time.Now().Unix()

	query := db.Database.Model(&entities.Paste{}).Where("expires < ?", currentTime)
	query.Delete(&entities.Paste{})
}
