package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {

	_, err := getUserByCookie(c)

	if err != nil {
		return c.Render(http.StatusOK, "auth.html", nil)
	}

	return c.Render(http.StatusOK, "home.html", nil)
}
