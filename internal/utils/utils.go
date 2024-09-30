package utils

import (
	"errors"
	"time"

	"github.com/robstave/rto/internal/domain/types"
)

// parseDate tries to parse a date string using multiple layouts.
// It returns the parsed time.Time or an error if none of the layouts match.
func ParseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",          // "YYYY-MM-DD"
		time.RFC3339,          // "YYYY-MM-DDTHH:MM:SSZ"
		"2006-01-02T15:04:05", // "YYYY-MM-DDTHH:MM:SS"
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("invalid date format")
}

// SameDay checks if two dates are on the same calendar day
func SameDay(a, b time.Time) bool {
	yearA, monthA, dayA := a.Date()
	yearB, monthB, dayB := b.Date()
	return yearA == yearB && monthA == monthB && dayA == dayB
}

// CalculateInOfficeAverage computes the number of in-office days and total days in the quarter
func CalculateInOfficeAverage(events []types.Event, startDate time.Time, endDate time.Time) (int, int) {
	// Define the quarter date range: October 1 to December 31 of the current year

	// Calculate total days in the quarter
	totalDays := int(endDate.Sub(startDate).Hours()/24) + 1 // +1 to include the end date

	inOfficeCount := 0

	// Iterate through all events and count in-office days within the quarter
	for _, event := range events {
		if event.Type == "attendance" && event.IsInOffice {
			if !event.Date.Before(startDate) && !event.Date.After(endDate) {
				inOfficeCount++
			}
		}
	}

	return inOfficeCount, totalDays
}

// GetCalendarMonth generates all weeks for the given month, including days from adjacent months
func GetCalendarMonth(currentDate time.Time) [][]types.CalendarDay {
	var weeks [][]types.CalendarDay

	// Normalize to the first day of the month
	firstOfMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())

	// Find the first Sunday before or on the first day of the month
	weekday := firstOfMonth.Weekday()
	daysToSubtract := int(weekday) // Sunday = 0
	startDate := firstOfMonth.AddDate(0, 0, -daysToSubtract)

	for week := 0; week < 6; week++ { // Up to 6 weeks in a month view
		var weekDays []types.CalendarDay
		for day := 0; day < 7; day++ {
			currentDay := startDate.AddDate(0, 0, week*7+day)
			inMonth := currentDay.Month() == firstOfMonth.Month()
			weekDays = append(weekDays, types.CalendarDay{
				Date:    currentDay,
				InMonth: inMonth,
			})
		}
		weeks = append(weeks, weekDays)

		// Check if all days in the current week are from the next month
		allDaysNextMonth := true
		for _, day := range weekDays {
			if day.Date.Month() == firstOfMonth.Month() {
				allDaysNextMonth = false
				break
			}
		}

		if allDaysNextMonth {
			weeks = weeks[:len(weeks)-1] // Remove the last week added
			break
		}
	}

	return weeks
}
