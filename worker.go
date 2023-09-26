package main

import (
	"time"

	"github.com/shair/handlers"
)

func runWorker() {
	for range time.Tick(time.Hour * 2) {
		handlers.DeleteExpiredPastes()
		handlers.DeleteExpiredUploads()
	}
}
