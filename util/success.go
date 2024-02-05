package util

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func CreateSuccessResult(c echo.Context, path string, token string) error {
	visitUrl := GetBaseUrl(c) + "/" + path + token

	data := echo.Map{
		"Token":      token,
		"VisitUrl":   visitUrl,
		"EncodedUrl": url.QueryEscape(visitUrl),
	}

	return c.Render(http.StatusCreated, "success.html", data)
}
