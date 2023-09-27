package config

import (
	"os"
	"strconv"

	"github.com/shair/util"
)

var TempLifeTime int64
var MaxFileSize int64
var UploadDir string
var DownloadDir string
var GalleryDir string
var AdminUsername string

func getInt64(envKey string, defaultValue int64) int64 {
	env := os.Getenv(envKey)
	val, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

func getString(envKey string, defaultValue string) string {
	env := os.Getenv(envKey)
	if util.IsBlank(env) {
		return defaultValue
	}
	return env
}

func Init() {
	TempLifeTime = getInt64("LIFETIME", 604800)         // default: 1 week
	MaxFileSize = getInt64("MAXFILESIZE", 1024*1024*50) // default: 50MB
	UploadDir = getString("UPLOADDIR", "./data/files")
	DownloadDir = getString("DOWNLOADDIR", "./data/downloads")
	GalleryDir = getString("GALLERYDIR", "./data/gallery")
	AdminUsername = getString("ADMINUSERNAME", "admin")

	for _, dir := range []string{UploadDir, DownloadDir, GalleryDir} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
		}
	}
}
