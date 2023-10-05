package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shair/config"
	"github.com/shair/db"
	"github.com/shair/entities"
	"github.com/shair/util"
)

const AuthCookie = "Auth"

func RegisterUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if util.IsBlank(username) || util.IsBlank(password) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Username and password can't be blank!",
		})
	}
	user, err := db.CreateUser(username, password)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"message": err.Error(),
		})
	}

	c.SetCookie(&http.Cookie{Name: AuthCookie, Value: user.AuthToken, Path: "/"})

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func LoginUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := db.LoginUser(username, password)

	if err == nil {
		c.SetCookie(&http.Cookie{Name: AuthCookie, Value: user.AuthToken, Path: "/"})
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func isAdmin(c echo.Context) bool {
	user, err := getUserByCookie(c)
	if err != nil || user.Username != config.AdminUsername {
		return false
	}
	return true
}

func getUserByCookie(c echo.Context) (entities.User, error) {
	var user entities.User

	authCookie, err := c.Cookie(AuthCookie)
	if err != nil || util.IsBlank(authCookie.Value) {
		return user, errors.New("Auth token missing!")
	}

	user, err = db.FindUser(authCookie.Value)
	if err != nil {
		return user, errors.New("User not found!")
	}

	return user, nil
}

func DeleteUser(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "User does not exist!",
		})
	}

	db.DeleteUser(user)

	return LogoutUser(c)
}

func LogoutUser(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: AuthCookie, Value: "", Path: "/"})

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
