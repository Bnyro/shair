package config

import (
	"os"
	"strconv"
)

var TempLifeTime int64
var MaxFileSize int64

func getInt64(envKey string, defaultValue int64) int64 {
	env := os.Getenv(envKey)
	val, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

func Init() {
	TempLifeTime = getInt64("LIFETIME", 604800)         // default: 1 week
	MaxFileSize = getInt64("MAXFILESIZE", 1024*1024*50) // default: 50MB
}
