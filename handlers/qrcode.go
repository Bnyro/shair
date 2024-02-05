package handlers

import (
	"image/png"
	"net/url"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/labstack/echo/v4"
	"github.com/shair/util"
)

func QrCode(c echo.Context) error {
	uri, err := url.QueryUnescape(c.QueryParam("url"))
	if err != nil {
		return err
	}

	parsedUri, err := url.Parse(uri)
	if err != nil || parsedUri.Host != util.GetHost(c) {
		return err
	}

	qrCode, err := qr.Encode(uri, qr.H, qr.Unicode)
	if err != nil {
		return err
	}

	image, err := barcode.Scale(qrCode, 300, 300)
	if err != nil {
		return err
	}

	return png.Encode(c.Response().Writer, image)
}
