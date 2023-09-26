package handlers

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/shair/config"
	"github.com/shair/util"
	"net/http"
	"os"
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

	return c.Render(http.StatusOK, "files.html", echo.Map{
		"Admin": user.Username == config.AdminUsername,
		"Files": util.WalkDir(config.DownloadDir),
	})
}

func DeleteFile(c echo.Context) error {
	if !isAdmin(c) {
		return errors.New("Not authorized!")
	}

	filename := c.Param("filename")
	os.Remove(config.DownloadDir + "/" + filename)

	return c.Redirect(http.StatusTemporaryRedirect, "/downloads/")
}
