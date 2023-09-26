package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/shair/config"
	"github.com/shair/util"
)

func Gallery(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return err
	}

	if file, err := c.FormFile("file"); err == nil {
		if file.Size > config.MaxFileSize {
			return errors.New("File size too big!")
		}
		err := util.Save(file, config.GalleryDir, file.Filename)
		if err != nil {
			fmt.Println(err)
		}
	}

	return c.Render(http.StatusOK, "gallery.html", echo.Map{
		"Admin":  user.Username == config.AdminUsername,
		"Images": util.WalkDir(config.GalleryDir),
	})
}

func DeleteImage(c echo.Context) error {
	if !isAdmin(c) {
		return errors.New("Not authorized!")
	}

	filename := c.Param("filename")
	os.Remove(config.GalleryDir + "/" + filename)

	return c.Redirect(http.StatusTemporaryRedirect, "/gallery/")
}
