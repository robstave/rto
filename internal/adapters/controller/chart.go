package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"
)

// GetChartData handles the retrieval of data for the D3 chart
func (ctlr *RTOController) GetChartData(c echo.Context) error {
	// Fetch all events from the service
	events := ctlr.service.GetAllEvents()

	// Build the date range (October 1, 2024 to December 31, 2024)
	startDate := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	dateRange := utils.GetDateRange(startDate, endDate) // We'll define this utility function next

	// Prepare data for the chart
	var data []map[string]interface{}

	total := 0.0
	for _, date := range dateRange {
		dayOfWeek := date.Weekday()
		isWeekday := dayOfWeek >= time.Monday && dayOfWeek <= time.Friday
		comesIn := isWeekday && (countAttendanceOnDate(events, date) > 0) // Define countAttendanceOnDate
		if comesIn {
			total += 7.0 / float64(len(dateRange))

		}

		data = append(data, map[string]interface{}{
			"date":    date.Format("2006-01-02"),
			"comesIn": comesIn,
			"total":   total,
		})
	}

	// Return the data as JSON
	return c.JSON(http.StatusOK, data)
}

// Utility function to count attendance on a specific date
func countAttendanceOnDate(events []types.Event, date time.Time) int {
	count := 0
	for _, event := range events {
		if utils.SameDay(event.Date, date) && event.Type == "attendance" && event.IsInOffice {
			count++
		}
	}
	return count
}
