package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/shair/entities"
)

func Save(file *multipart.FileHeader, dir string, name string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dir + "/" + name)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func WalkDir(path string) []entities.File {
	var files []entities.File

	content, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		return files
	}
	for _, file := range content {
		if file.IsDir() {
			continue
		}
		stats, _ := os.Stat(path + "/" + file.Name())

		files = append(files, entities.File{
			Name:    file.Name(),
			Size:    humanize.Bytes(uint64(stats.Size())),
			ModTime: stats.ModTime().Format(time.RFC822),
		})
	}

	return files
}
