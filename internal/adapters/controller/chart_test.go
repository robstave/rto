// controller/get_chart_data_test.go

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

var (
	QuarterStart = time.Date(2024, time.October, 30, 0, 0, 0, 0, time.UTC)
	QuarterEnd   = time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC)
)

func TestGetChartData_Success(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define mock preferences
	prefs := types.Preferences{
		TargetDays: "2.5",
	}

	// Define mock events
	events := []types.Event{
		{
			ID:          1,
			Date:        time.Date(2024, 10, 2, 0, 0, 0, 0, time.UTC),
			Description: "Attendance",
			Type:        "attendance",
			IsInOffice:  true,
		},
		{
			ID:          2,
			Date:        time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC),
			Description: "Remote",
			Type:        "attendance",
			IsInOffice:  false,
		},
	}

	// Setup expectations
	mockService.On("GetAllEvents").Return(events)
	mockService.On("GetPrefs").Return(prefs)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService, QuarterStart, QuarterEnd)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a GET request
	req := httptest.NewRequest(http.MethodGet, "/chart-data", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.GetChartData(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)

		// Optionally, you can unmarshal the response and verify its structure
		var response ChartResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, prefs.TargetDays, fmt.Sprintf("%.1f", response.TargetDays))
		assert.Len(t, response.Data, 92) // From Oct 1 to Dec 31 is 92 days
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}

func TestGetChartData_TargetDaysParsingError(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define mock preferences with invalid TargetDays
	prefs := types.Preferences{
		TargetDays: "invalid",
	}

	// Define mock events
	events := []types.Event{}

	// Setup expectations
	mockService.On("GetAllEvents").Return(events)
	mockService.On("GetPrefs").Return(prefs)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService, QuarterStart, QuarterEnd)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a GET request
	req := httptest.NewRequest(http.MethodGet, "/chart-data", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.GetChartData(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)

		// Optionally, you can unmarshal the response and verify its structure
		var response ChartResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2.5, response.TargetDays) // Fallback to default
		assert.Len(t, response.Data, 92)          // From Oct 1 to Dec 31 is 92 days
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}

func TestGetChartData_EmptyEvents(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define mock preferences
	prefs := types.Preferences{
		TargetDays: "3.0",
	}

	// Define empty mock events
	events := []types.Event{}

	// Setup expectations
	mockService.On("GetAllEvents").Return(events)
	mockService.On("GetPrefs").Return(prefs)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService, QuarterStart, QuarterEnd)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a GET request
	req := httptest.NewRequest(http.MethodGet, "/chart-data", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.GetChartData(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)

		// Optionally, you can unmarshal the response and verify its structure
		var response ChartResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response.Data, 92) // From Oct 1 to Dec 31 is 92 days
		// Further checks can be made on the content of 'data'
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}
