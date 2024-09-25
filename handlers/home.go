package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// CalendarDay represents a single day in the calendar
type CalendarDay struct {
	Date    time.Time
	InMonth bool
}

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

		// Stop if we've passed the end of the month and the week contains only days from the next month
		if weekDays[0].Date.Month() > firstOfMonth.Month() || (weekDays[0].Date.Month() == firstOfMonth.Month() && weekDays[0].Date.Day() == 1 && week > 0) {
			break
		}
	}

	return weeks
}
