// Package helpers provides different small helpers
package helpers

import (
	"time"
)

// TimeNow returns current UTC time
func TimeNow() time.Time {
	return time.Now().UTC()
}

// TimeNowAddHours returns current UTC time extend with given amount of hours
func TimeNowAddHours(hours int) time.Time {
	return TimeNow().Add(time.Duration(hours) * time.Hour)
}
