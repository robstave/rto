// controller/add_event_test.go

package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/internal/domain/mocks"
	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddEvent_Success(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define the event to be added
	event := types.Event{
		Date:        time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC),
		Description: "Team Building",
		Type:        "vacation",
		IsInOffice:  false,
	}

	// Setup expectations
	mockService.On("AddEvent", event).Return(nil)
	mockService.On("GetEventByDateAndType", event.Date, "vacation").Return(&event, nil)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	reqBody := "date=2024-11-20&type=vacation&description=Team+Building&isInOffice=false"

	req := httptest.NewRequest(http.MethodPost, "/add-event", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mockRenderer{}

	// Call the handler
	if assert.NoError(t, ctlr.AddEvent(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)
		// You can add more assertions based on your handler's behavior, such as checking the rendered template
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}

func TestAddEvent_MissingFields(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	reqBody := "date=2024-11-20&description=Team+Building&isInOffice=false"

	req := httptest.NewRequest(http.MethodPost, "/add-event", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.AddEvent(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponse := `{"success": false, "message": "Date and Event Type are required"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

}

func TestAddEvent_InvalidDateFormat(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request with invalid date format
	reqBody := "date=11-2024-20&type=vacation&description=Team+Building&isInOffice=false"

	req := httptest.NewRequest(http.MethodPost, "/add-event", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.AddEvent(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponse := `{"success": false, "message": "Invalid date format"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Ensure that the service was not called
	mockService.AssertNotCalled(t, "AddEvent", mock.Anything)
}

func TestAddEvent_ServiceError(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define the event to be added
	event := types.Event{
		Date:        time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC),
		Description: "Team Building",
		Type:        "vacation",
		IsInOffice:  false,
	}

	// Setup expectations
	mockService.On("AddEvent", event).Return(errors.New("database error"))

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request with the event data
	reqBody := "date=2024-11-20&type=vacation&description=Team+Building&isInOffice=false"

	req := httptest.NewRequest(http.MethodPost, "/add-event", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.AddEvent(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		expectedResponse := `{"success": false, "message": "Failed to add event."}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}

func TestBulkAddEventsJSON_Success(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define mock bulk add response
	bulkAddResponse := types.BulkAddResponse{
		Success: true,
		Added:   2,
		Updated: 1,
		Skipped: 0,
		Message: "Successfully added 2 vacation events and updated 1 event.",
		Results: []types.BulkAddResult{
			{
				Date:        "2024-11-20",
				Action:      "Added new vacation",
				Description: "Conference",
			},
			{
				Date:        "2024-11-21",
				Action:      "Added new vacation",
				Description: "Workshop",
			},
			{
				Date:        "2024-11-22",
				Action:      "Updated existing vacation",
				Description: "Team Meeting",
			},
		},
	}

	// Setup expectations
	events := []types.Event{
		{
			Date:        time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC),
			Description: "Conference",
			Type:        "vacation",
		},
		{
			Date:        time.Date(2024, 11, 21, 0, 0, 0, 0, time.UTC),
			Description: "Workshop",
			Type:        "vacation",
		},
		{
			Date:        time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC),
			Description: "Team Meeting",
			Type:        "vacation",
		},
	}
	mockService.On("BulkAddEvents", events).Return(&bulkAddResponse, nil)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request with bulk events
	reqBody := BulkAddEventsJSONRequest{
		Events: []RawEvent{
			{
				Date:        "2024-11-20",
				Description: "Conference",
				Type:        "vacation",
			},
			{
				Date:        "2024-11-21",
				Description: "Workshop",
				Type:        "vacation",
			},
			{
				Date:        "2024-11-22",
				Description: "Team Meeting",
				Type:        "vacation",
			},
		},
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/add-events-json", bytes.NewReader(reqBodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.BulkAddEventsJSON(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse, _ := json.Marshal(bulkAddResponse)
		assert.JSONEq(t, string(expectedResponse), rec.Body.String())
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}

func TestBulkAddEventsJSON_InvalidJSON(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request with invalid JSON
	invalidJSON := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/add-events-json", bytes.NewBufferString(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.BulkAddEventsJSON(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponse := `{"success": false, "message": "Invalid JSON payload."}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Ensure that the service was not called
	mockService.AssertNotCalled(t, "BulkAddEvents", mock.Anything)
}

func TestBulkAddEventsJSON_InvalidEventType(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request with an invalid event type
	reqBody := BulkAddEventsJSONRequest{
		Events: []RawEvent{
			{
				Date:        "2024-11-20",
				Description: "Conference",
				Type:        "invalidType", // Invalid type
			},
		},
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/add-events-json", bytes.NewReader(reqBodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.BulkAddEventsJSON(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponse := `{"success": false, "message": "Only events with type 'vacation' can be bulk added (error at index 0)."}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Ensure that the service was not called
	mockService.AssertNotCalled(t, "BulkAddEvents", mock.Anything)
}

func TestBulkAddEventsJSON_InvalidDateFormat(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request with an invalid date format
	reqBody := BulkAddEventsJSONRequest{
		Events: []RawEvent{
			{
				Date:        "20-11-2024", // Invalid format
				Description: "Conference",
				Type:        "vacation",
			},
		},
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/add-events-json", bytes.NewReader(reqBodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.BulkAddEventsJSON(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponse := `{"success": false, "message": "Invalid date format for event on 20-11-2024. Expected YYYY-MM-DD."}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Ensure that the service was not called
	mockService.AssertNotCalled(t, "BulkAddEvents", mock.Anything)
}

func TestBulkAddEventsJSON_ServiceError(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Define mock bulk add response with an error
	mockService.On("BulkAddEvents", mock.Anything).Return(nil, errors.New("service error"))

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a POST request with bulk events
	reqBody := BulkAddEventsJSONRequest{
		Events: []RawEvent{
			{
				Date:        "2024-11-20",
				Description: "Conference",
				Type:        "vacation",
			},
		},
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/add-events-json", bytes.NewReader(reqBodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if assert.NoError(t, ctlr.BulkAddEventsJSON(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		expectedResponse := `{"success": false, "message": "An error occurred while processing bulk add events."}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}
