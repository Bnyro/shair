package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"github.com/shair/config"
	"github.com/shair/db"
	"github.com/shair/entities"
	"github.com/shair/util"
)

func NewUpload(c echo.Context) error {
	if c.Request().Method == "GET" {
		return c.Render(http.StatusOK, "newupload.html", nil)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	if file.Size > config.MaxFileSize {
		return errors.New("File too big!")
	}

	token := util.GenerateSecureToken(30)

	if err = util.Save(file, config.UploadDir, token); err != nil {
		return err
	}

	now := time.Now().Unix()
	upload := entities.Upload{
		Name:    file.Filename,
		Size:    file.Size,
		Created: now,
		Expires: now + config.TempLifeTime,
		Token:   token,
	}

	db.Database.Create(&upload)

	return util.CreateSuccessResult(c, "upload", upload.Token)
}

func GetUpload(c echo.Context) error {
	upload := entities.Upload{
		Token: c.Param("token"),
	}

	if db.Database.Where(&upload).Find(&upload).RowsAffected == 0 {
		c.Render(http.StatusNotFound, "notfound.html", nil)
	}

	data := echo.Map{
		"Upload": upload,
		"Size":   humanize.Bytes(uint64(upload.Size)),
	}

	return c.Render(http.StatusOK, "upload.html", data)
}

func DeleteExpiredUploads() {
	var uploads []entities.Upload

	currentTime := time.Now().Unix()
	query := db.Database.Model(&entities.Upload{}).Where("expires < ?", currentTime)

	query.Find(&uploads)

	for _, upload := range uploads {
		os.Remove(config.UploadDir + "/" + upload.Token)
	}

	query.Delete(&uploads)
}
