package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"sync"

	"github.com/labstack/echo/v4"
)

// CalendarDay represents a single day in the calendar
type CalendarDay struct {
	Date    time.Time
	InMonth bool
	Today   bool
	Events  []Event
}

// Assume events are now a global variable or part of a struct for persistence
//var allEvents []Event

// Global variable to store all events and manage thread safety
var (
	allEvents  []Event
	eventsLock sync.RWMutex
)

// Home renders the calendar on the home page
func Home(c echo.Context) error {
	c.Logger().Info("------")
	c.Logger().Info("------")
	c.Logger().Info("------")

	// Get current date or date from query parameters
	currentDate := time.Now()
	yearParam := c.QueryParam("year")
	monthParam := c.QueryParam("month")
	dayParam := c.QueryParam("day")

	if yearParam != "" && monthParam != "" && dayParam != "" {
		year, err1 := strconv.Atoi(yearParam)
		month, err2 := strconv.Atoi(monthParam)
		day, err3 := strconv.Atoi(dayParam)
		if err1 == nil && err2 == nil && err3 == nil {
			currentDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, currentDate.Location())
		}
	}

	// Generate calendar for the current month
	weeks := getCalendarMonth(currentDate)

	// Precompute formatted dates for navigation links
	prevMonthDate := currentDate.AddDate(0, -1, 0)
	nextMonthDate := currentDate.AddDate(0, 1, 0)

	// Assign events to the corresponding days
	for weekIdx, week := range weeks {
		for dayIdx, day := range week {
			dateStr := day.Date.Format("2006-01-02") // YYYY-MM-DD
			dayEvents := []Event{}

			for _, event := range allEvents {
				if event.Date.Format("2006-01-02") == dateStr {
					dayEvents = append(dayEvents, event)
				}
			}

			//// Convert to event descriptions for the template
			//var eventDescriptions []string
			//for _, e := range dayEvents {
			//	eventDescriptions = append(eventDescriptions, e.Description)
			//}

			weeks[weekIdx][dayIdx].Events = dayEvents

			// Mark today if applicable
			//if day.Date.Equal(today) {
			//	weeks[weekIdx][dayIdx].Today = true
			//}
		}
	}

	// Define 'today' before the loop
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// Assign Today Flag
	for weekIdx, week := range weeks {
		for dayIdx, day := range week {
			if day.Date.Equal(today) {
				weeks[weekIdx][dayIdx].Today = true
			}
		}
	}

	data := map[string]interface{}{
		"CurrentDate": currentDate,
		"Weeks":       weeks,
		"PrevMonth": map[string]string{
			"year":  prevMonthDate.Format("2006"),
			"month": prevMonthDate.Format("01"),
			"day":   prevMonthDate.Format("02"),
		},
		"NextMonth": map[string]string{
			"year":  nextMonthDate.Format("2006"),
			"month": nextMonthDate.Format("01"),
			"day":   nextMonthDate.Format("02"),
		},
	}

	//log
	for weekIdx, week := range weeks {
		for dayIdx, day := range week {
			// Existing event assignment logic

			// Debugging: Log events for each day
			if len(weeks[weekIdx][dayIdx].Events) > 0 {
				log.Printf("Date: %s, Events: %+v\n", day.Date.Format("2006-01-02"), weeks[weekIdx][dayIdx].Events)
			}
		}
	}

	// Render the template
	if err := c.Render(http.StatusOK, "home.html", data); err != nil {
		c.Logger().Error("Template rendering error:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return nil
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

	// Iterate over the weeks
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
