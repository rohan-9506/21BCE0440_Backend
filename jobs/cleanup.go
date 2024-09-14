package jobs

import (
	"time"
)

func StartCleanupJob() {
	for {
		// Perform file cleanup tasks here
		time.Sleep(24 * time.Hour) // Run cleanup daily
	}
}
