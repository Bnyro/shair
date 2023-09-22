package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"github.com/shair/config"
	"github.com/shair/entities"
	"github.com/shair/util"
)

func Files(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return err
	}

	if file, err := c.FormFile("file"); err == nil {
		if file.Size > config.MaxFileSize {
			return errors.New("File size too big!")
		}
		err := util.Save(file, config.DownloadDir, file.Filename)
		if err != nil {
			fmt.Println(err)
		}
	}

	content, err := os.ReadDir(config.DownloadDir)
	if err != nil {
		fmt.Println(err.Error())
	}

	var files []entities.File
	for _, file := range content {
		if file.IsDir() {
			continue
		}

		stats, _ := os.Stat(config.DownloadDir + "/" + file.Name())
		files = append(files, entities.File{
			Name:    file.Name(),
			Size:    humanize.Bytes(uint64(stats.Size())),
			ModTime: stats.ModTime().Format(time.RFC822),
		})
	}

	return c.Render(http.StatusOK, "files.html", echo.Map{
		"Admin": user.Username == config.AdminUsername,
		"Files": files,
	})
}

func DeleteFile(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil || user.Username != config.AdminUsername {
		return err
	}

	filename := c.Param("filename")
	os.Remove(config.DownloadDir + "/" + filename)

	return c.Redirect(http.StatusTemporaryRedirect, "/downloads/")
}
