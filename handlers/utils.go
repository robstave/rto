package handlers

import (
	"time"
)

// sameDay checks if two dates are on the same calendar day
func sameDay(a, b time.Time) bool {
	yearA, monthA, dayA := a.Date()
	yearB, monthB, dayB := b.Date()
	return yearA == yearB && monthA == monthB && dayA == dayB
}

// calculateInOfficeAverage computes the number of in-office days and total days in the quarter
func calculateInOfficeAverage(events []Event, startDate time.Time, endDate time.Time) (int, int) {
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

// getCalendarMonth generates all weeks for the given month, including days from adjacent months
func getCalendarMonth(currentDate time.Time) [][]CalendarDay {
	var weeks [][]CalendarDay

	// Normalize to the first day of the month
	firstOfMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())

	// Find the first Sunday before or on the first day of the month
	weekday := firstOfMonth.Weekday()
	daysToSubtract := int(weekday) // Sunday = 0
	startDate := firstOfMonth.AddDate(0, 0, -daysToSubtract)

	for week := 0; week < 6; week++ { // Up to 6 weeks in a month view
		var weekDays []CalendarDay
		for day := 0; day < 7; day++ {
			currentDay := startDate.AddDate(0, 0, week*7+day)
			inMonth := currentDay.Month() == firstOfMonth.Month()
			weekDays = append(weekDays, CalendarDay{
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
