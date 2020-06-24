package utils

import (
	"time"
)

func GetLocation(timezone string) *time.Location {
	loc, err := time.LoadLocation(timezone)
	if err == nil {
		return loc
	}
	return time.UTC
}
