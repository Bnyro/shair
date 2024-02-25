package util

import "github.com/labstack/echo/v4"

func GetHost(c echo.Context) string {
	// if behind a proxy, it is necessary to set the "X-Forwarded-Host" header to run properly
	forwardedFor := c.Request().Header.Get("X-Forwarded-Host")

	if !IsBlank(forwardedFor) {
		return forwardedFor
	}

	return c.Request().Host
}

func GetScheme(c echo.Context) string {
	scheme := c.Request().Header.Get("X-Forwarded-Proto")

	if IsBlank(scheme) {
		scheme = c.Request().URL.Scheme
	}

	return scheme
}

func GetBaseUrl(c echo.Context) string {
	return GetScheme(c) + "://" + GetHost(c)
}
