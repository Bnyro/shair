package util

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func IsBlank(text string) bool {
	return len(strings.TrimSpace(text)) == 0
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GetBaseUrl(c echo.Context) string {
	scheme := c.Request().URL.Scheme
	if IsBlank(scheme) {
		scheme = "http"
	}
	return scheme + "://" + c.Request().Host
}
