package types

import (
	"time"
)

type Event struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	IsInOffice  bool      `json:"isInOffice"`
	Type        string    `json:"type"` // "attendance", "holiday", "vacation"
}

type Preferences struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	DefaultDays string `json:"defaultDays"` // e.g., "M,T,W,Th,F"
	TargetDays  string `json:"targetDays"`  // e.g., "2.5"
}

// CalendarDay represents a single day in the calendar
type CalendarDay struct {
	Date      time.Time
	InMonth   bool
	Today     bool
	Events    []Event
	IsWeekend bool // New field to indicate weekends

}

type AttendanceStats struct {
	InOfficeCount int
	TotalDays     int
	Average       float64
	AverageDays   float64
	TargetDays    float64
}
