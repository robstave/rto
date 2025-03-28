package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"
)

// ChartResponse represents the JSON response for the chart data
type ChartResponse struct {
	Data       []map[string]interface{} `json:"data"`
	TargetDays float64                  `json:"targetDays"`
}

// GetChartData handles the retrieval of data for the D3 chart
func (ctlr *RTOController) GetChartData(c echo.Context) error {
	// Fetch all events from the service
	events := ctlr.service.GetAllEvents()
	prefs := ctlr.service.GetPrefs()
	targetDaysStr := prefs.TargetDays
	targetDays, err := strconv.ParseFloat(targetDaysStr, 64)
	if err != nil {
		// Fallback to default target if parsing fails
		targetDays = 2.5
	}

	// Build the date range (October 1, 2024 to December 31, 2024)
	//startDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	//endDate := time.Date(2025, 3, 28, 0, 0, 0, 0, time.UTC)
	startDate := ctlr.quarterStart
	endDate := ctlr.quarterEnd
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

	// Create the ChartResponse
	response := ChartResponse{
		Data:       data,
		TargetDays: targetDays,
	}

	// Return the ChartResponse as JSON
	return c.JSON(http.StatusOK, response)

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
