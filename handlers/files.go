package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/shair/config"
	"github.com/shair/util"
	"net/http"
	"os"
)

func Files(c echo.Context) error {
	isUserAdmin := isAdmin(c)
	if file, err := c.FormFile("file"); err == nil && isUserAdmin {
		util.Save(file, config.DownloadDir, file.Filename)
	}

	return c.Render(http.StatusOK, "files.html", echo.Map{
		"Admin": isUserAdmin,
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
