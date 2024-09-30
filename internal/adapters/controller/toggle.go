package controller

import (
	"net/http"

	"time"

	"github.com/labstack/echo/v4"
)

// ToggleAttendanceRequest represents the JSON payload for toggling attendance
type ToggleAttendanceRequest struct {
	Date string `json:"date"` // Expected format: YYYY-MM-DD
}

// ToggleAttendanceResponse represents the JSON response after toggling
type ToggleAttendanceResponse struct {
	Success       bool    `json:"success"`
	NewStatus     string  `json:"newStatus,omitempty"` // "in" or "remote"
	Message       string  `json:"message,omitempty"`
	InOfficeCount int     `json:"inOfficeCount,omitempty"`
	TotalDays     int     `json:"totalDays,omitempty"`
	Average       float64 `json:"average,omitempty"`
	AverageDays   float64 `json:"averageDays,omitempty"`
	TargetDays    float64 `json:"targetDays,omitempty"` // New field for target value
}

// ToggleAttendance handles toggling attendance status for a given date
func (ctlr *RTOController) ToggleAttendance(c echo.Context) error {
	req := new(ToggleAttendanceRequest)
	if err := c.Bind(req); err != nil {
		ctlr.logger.Error("Error binding request", "error", err)
		return c.JSON(http.StatusBadRequest, ToggleAttendanceResponse{
			Success: false,
			Message: "Invalid request payload.",
		})
	}

	// Parse the date
	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		ctlr.logger.Error("Error parsing date", "fn", "ToggleAttendance", "error", err)
		return c.JSON(http.StatusBadRequest, ToggleAttendanceResponse{
			Success: false,
			Message: "Invalid date format. Expected YYYY-MM-DD.",
		})
	}

	// Call domain service to toggle attendance
	newStatus, err := ctlr.service.ToggleAttendance(eventDate)
	if err != nil {
		ctlr.logger.Error("Error toggling attendance", "error", err)
		return c.JSON(http.StatusInternalServerError, ToggleAttendanceResponse{
			Success: false,
			Message: "Failed to toggle attendance.",
		})
	}

	// After toggling, recalculate stats
	stats, err := ctlr.service.CalculateAttendanceStats()
	if err != nil {
		ctlr.logger.Error("Error calculating stats", "error", err)
		return c.JSON(http.StatusInternalServerError, ToggleAttendanceResponse{
			Success: false,
			Message: "Failed to calculate attendance statistics.",
		})
	}

	return c.JSON(http.StatusOK, ToggleAttendanceResponse{
		Success:       true,
		NewStatus:     newStatus,
		InOfficeCount: stats.InOfficeCount,
		TotalDays:     stats.TotalDays,
		Average:       stats.Average,
		AverageDays:   stats.AverageDays,
		TargetDays:    stats.TargetDays,
	})
}
