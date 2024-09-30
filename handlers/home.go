package handlers

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"

	"github.com/labstack/echo/v4"
)

// Global variable to store all events and manage thread safety
var (
	allEvents  []types.Event
	eventsLock sync.RWMutex
)

// Home renders the calendar on the home page
func Home(c echo.Context) error {

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
	weeks := utils.GetCalendarMonth(currentDate)

	// Precompute formatted dates for navigation links
	prevMonthDate := currentDate.AddDate(0, -1, 0)
	nextMonthDate := currentDate.AddDate(0, 1, 0)

	// Assign events to the corresponding days
	for weekIdx, week := range weeks {
		for dayIdx, day := range week {
			dateStr := day.Date.Format("2006-01-02") // YYYY-MM-DD
			dayEvents := []types.Event{}

			for _, event := range allEvents {
				if event.Date.Format("2006-01-02") == dateStr {
					dayEvents = append(dayEvents, event)
				}
			}

			weeks[weekIdx][dayIdx].Events = dayEvents

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
	currentYear := time.Now().Year()
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.Local)
	// Calculate In-Office Average
	inOfficeCount, totalDays := utils.CalculateInOfficeAverage(allEvents, startDate, endDate)

	average := 0.0
	averageDays := 0.0
	if totalDays > 0 {
		average = (float64(inOfficeCount) / float64(totalDays)) * 100
		averageDays = (float64(inOfficeCount) / float64(totalDays)) * 7 //average days/week
	}

	// Fetch target days from preferences
	preferencesLock.RLock()
	currentPreferences := preferences
	preferencesLock.RUnlock()

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
		"InOfficeCount": inOfficeCount,
		"TotalDays":     totalDays,
		"Average":       average,
		"AverageDays":   averageDays,
		"Preferences":   currentPreferences, // Add Preferences here

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
		logger.Error("Template rendering error:", "error", err)

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return nil
}
