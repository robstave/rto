package handlers

import (
	"time"
)

// isWeekend determines if a given date is Saturday or Sunday
func isWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}
