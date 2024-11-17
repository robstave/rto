// controller/toggle_attendance_test.go

package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/internal/domain/mocks"
	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestToggleAttendance_Success(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Set up request payload
	reqBody := ToggleAttendanceRequest{
		Date: "2024-10-15",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	// Mock the service methods
	eventDate, _ := time.Parse("2006-01-02", reqBody.Date)
	mockService.On("ToggleAttendance", eventDate).Return("in", nil)
	mockService.On("CalculateAttendanceStats").Return(&types.AttendanceStats{
		InOfficeCount:  10,
		TotalDays:      20,
		Average:        50.0,
		AverageDays:    3.5,
		TargetDays:     2.5,
		AveragePercent: 140.0,
	}, nil)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request
	req := httptest.NewRequest(http.MethodPost, "/toggle-attendance", bytes.NewReader(reqBodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.ToggleAttendance(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse := `{
			"success": true,
			"newStatus": "in",
			"inOfficeCount": 10,
			"totalDays": 20,
			"average": 50.0,
			"averageDays": 3.5,
			"targetDays": 2.5
		}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}

func TestToggleAttendance_InvalidDateFormat(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Set up request payload with invalid date format
	reqBody := ToggleAttendanceRequest{
		Date: "15-10-2024", // Invalid format
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request
	req := httptest.NewRequest(http.MethodPost, "/toggle-attendance", bytes.NewReader(reqBodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.ToggleAttendance(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponse := `{
			"success": false,
			"message": "Invalid date format. Expected YYYY-MM-DD."
		}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestToggleAttendance_ServiceError(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Set up request payload
	reqBody := ToggleAttendanceRequest{
		Date: "2024-10-15",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	// Mock the service method to return an error
	eventDate, _ := time.Parse("2006-01-02", reqBody.Date)
	mockService.On("ToggleAttendance", eventDate).Return("", errors.New("attendance event not found"))

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request
	req := httptest.NewRequest(http.MethodPost, "/toggle-attendance", bytes.NewReader(reqBodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.ToggleAttendance(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		expectedResponse := `{
			"success": false,
			"message": "Failed to toggle attendance."
		}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}
