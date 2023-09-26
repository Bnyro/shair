package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/shair/config"
	"github.com/shair/util"
)

func Gallery(c echo.Context) error {
	isUserAdmin := isAdmin(c)
	if file, err := c.FormFile("file"); err == nil && isUserAdmin {
		util.Save(file, config.GalleryDir, file.Filename)
	}

	return c.Render(http.StatusOK, "gallery.html", echo.Map{
		"Admin":  isUserAdmin,
		"Images": util.WalkDir(config.GalleryDir),
	})
}

func GalleryDia(c echo.Context) error {
	return c.Render(http.StatusOK, "dia.html", echo.Map{
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
