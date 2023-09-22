package util

import (
	"io"
	"mime/multipart"
	"os"
)

func Save(file *multipart.FileHeader, token string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create("files/" + token)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
