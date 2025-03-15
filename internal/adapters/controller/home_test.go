// controller/home_test.go

package controller

import (
	"encoding/json"
	"fmt"
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

func TestHome_Success(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define mock events
	mockEvents := []types.Event{
		{
			ID:          1,
			Date:        time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC),
			Description: "Conference",
			Type:        "vacation",
			IsInOffice:  false,
		},
	}

	// Define mock preferences
	prefs := types.Preferences{
		TargetDays: "2.5",
	}

	// Define mock attendance stats
	attendanceStats := &types.AttendanceStats{
		InOfficeCount:  0,
		TotalDays:      92,
		Average:        25.0,
		AverageDays:    1.75,
		TargetDays:     2.5,
		AveragePercent: 70.0,
	}

	// Setup expectations
	mockService.On("GetAllEvents").Return(mockEvents)
	mockService.On("GetPrefs").Return(prefs)
	//mockService.On("CalculateAttendanceStats").Return(attendanceStats, nil)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService, QuarterStart, QuarterEnd)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a GET request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assign a mock renderer
	e.Renderer = &mockRenderer{}

	// Call the handler
	if assert.NoError(t, ctlr.Home(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)

		// Optionally, verify the rendered content
		// Since we're using a mock renderer that serializes data to JSON, we can unmarshal and verify
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, prefs.TargetDays, fmt.Sprintf("%.1f", response["TargetDays"].(float64)))
		assert.Equal(t, attendanceStats.InOfficeCount, int(response["InOfficeCount"].(float64)))
		assert.Equal(t, attendanceStats.TotalDays, int(response["TotalDays"].(float64)))
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}

func TestHome_NoEvents(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define mock preferences
	prefs := types.Preferences{
		TargetDays: "2.5",
	}

	// Define mock attendance stats
	attendanceStats := &types.AttendanceStats{
		InOfficeCount:  0,
		TotalDays:      92, // October to December
		Average:        0.0,
		AverageDays:    0.0,
		TargetDays:     2.5,
		AveragePercent: 0.0,
	}

	// Setup expectations
	mockService.On("GetAllEvents").Return([]types.Event{})
	mockService.On("GetPrefs").Return(prefs)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService, QuarterStart, QuarterEnd)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a GET request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assign a mock renderer
	e.Renderer = &mockRenderer{}

	// Call the handler
	if assert.NoError(t, ctlr.Home(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)

		// Optionally, verify the rendered content
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, attendanceStats.InOfficeCount, int(response["InOfficeCount"].(float64)))
		assert.Equal(t, attendanceStats.TotalDays, int(response["TotalDays"].(float64)))
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}
