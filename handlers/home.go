package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {

	_, err := getUserByCookie(c)

	if err != nil {
		return c.Render(http.StatusOK, "auth.html", nil)
	}

	return c.Render(http.StatusOK, "home.html", nil)
}

func SetTheme(c echo.Context) error {
	theme := c.FormValue("theme")
	c.SetCookie(&http.Cookie{Name: "Theme", Value: theme, Path: "/"})

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Status(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}
