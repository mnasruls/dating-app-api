package helpers

import "time"

func GetTimeToMidnight() int {
	// Get the current time
	now := time.Now()

	// Get the start of the next day (midnight)
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	// Calculate the duration until next midnight
	durationUntilMidnight := nextMidnight.Sub(now)

	return int(durationUntilMidnight.Hours())
}
