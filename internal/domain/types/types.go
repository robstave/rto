package types

import (
	"fmt"
	"time"
)

type Event struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Date        time.Time `gorm:"type:date;not null"` // Use 'date' type to store only the date
	Description string    `gorm:"type:varchar(255);not null"`
	Type        string    `gorm:"type:varchar(50);not null"` // "holiday", "vacation", "attendance"
	IsInOffice  bool      `gorm:"default:false"`             // Relevant for "attendance" type
}

func (e Event) String() string {
	return fmt.Sprintf("Event{ID: %d, Date: %s, Description: %q, IsInOffice: %t, Type: %q}",
		e.ID,
		e.Date.Format("2006-01-02"),
		e.Description,
		e.IsInOffice,
		e.Type)
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
	IsWeekend bool //
	IsFuture  bool // future days

}

type AttendanceStats struct {
	InOfficeCount  int
	TotalDays      int
	Average        float64
	AverageDays    float64
	TargetDays     float64
	AveragePercent float64
}

type BulkAddResult struct {
	Date        string `json:"date"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
	Error       string `json:"error,omitempty"`
}

// BulkAddResponse encapsulates the overall result of a bulk add operation
type BulkAddResponse struct {
	Success bool            `json:"success"`
	Added   int             `json:"added"`
	Updated int             `json:"updated"`
	Skipped int             `json:"skipped"`
	Message string          `json:"message"`
	Results []BulkAddResult `json:"results"`
}
