package types

import (
	"time"
)

// Event represents a calendar event
type Event struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	IsInOffice  bool      `json:"isInOffice"`
	Type        string    `json:"type"` // "attendance", "holiday", "vacation"
}

// CalendarDay represents a single day in the calendar
type CalendarDay struct {
	Date      time.Time
	InMonth   bool
	Today     bool
	Events    []Event
	IsWeekend bool // New field to indicate weekends

}
