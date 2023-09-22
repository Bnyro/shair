package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateSuccessResult(c echo.Context, path string, token string) error {
	data := echo.Map{
		"Token":   token,
		"Path":    path,
		"BaseUrl": GetBaseUrl(c),
	}
	return c.Render(http.StatusCreated, "success.html", data)
}
