package handlers

import (
	"errors"
	"image/png"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/shair/config"
	"github.com/shair/db"
	"github.com/shair/entities"
	"github.com/shair/util"
)

const AuthCookie = "Auth"

func Auth(c echo.Context) error {
	return c.Render(http.StatusOK, "auth.html", nil)
}

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
	totp := c.FormValue("totp")

	user, err := db.LoginUser(username, password, totp)

	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth")
	}

	c.SetCookie(&http.Cookie{Name: AuthCookie, Value: user.AuthToken, Path: "/"})
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

func SetupTotp(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return err
	}

	totpSecret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Shair",
		AccountName: user.Username,
	})

	if err != nil {
		return err
	}

	qrCodeSource := url.QueryEscape(totpSecret.URL())
	return c.Render(http.StatusOK, "totp.html", echo.Map{
		"Secret":    totpSecret.Secret(),
		"QrCodeUrl": qrCodeSource,
	})
}

func ValidateTotp(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return err
	}

	secret := c.FormValue("secret")

	validation, _ := url.QueryUnescape(c.FormValue("validation"))
	if err != nil {
		return err
	}

	if !totp.Validate(validation, secret) {
		return c.Redirect(http.StatusTemporaryRedirect, "/totp")
	}

	user.TotpSecret = secret
	db.Database.Where("id = ?", user.ID).Updates(&user)

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func TotpQrCode(c echo.Context) error {
	otpauthUrl, err := url.QueryUnescape(c.QueryParam("url"))
	if err != nil {
		return err
	}

	key, err := otp.NewKeyFromURL(otpauthUrl)
	if err != nil {
		return err
	}

	image, err := key.Image(300, 300)
	if err != nil {
		return err
	}

	return png.Encode(c.Response().Writer, image)
}

func DisableTotp(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return err
	}

	db.Database.Model(&user).Where("id = ?", user.ID).Update("totp_secret", "")

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
