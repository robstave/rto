// controller/delete_event_test.go

package controller

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEvent_Success(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL (Business Logic Layer)
	mockService := new(mocks.RTOBLL)

	// Mock the service method to return no error
	mockService.On("TransformVacationToRemote", 1).Return(nil)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService, QuarterStart, QuarterEnd)

	// Create a DELETE request with a valid event ID
	req := httptest.NewRequest(http.MethodDelete, "/events/delete/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the handler
	if assert.NoError(t, ctlr.DeleteEvent(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse := `{"success": true, "message": "Vacation day transformed into a remote day successfully."}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}

func TestDeleteEvent_InvalidID(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService, QuarterStart, QuarterEnd)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger
	// Create a DELETE request with an invalid event ID
	req := httptest.NewRequest(http.MethodDelete, "/events/delete/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	// Call the handler
	if assert.NoError(t, ctlr.DeleteEvent(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponse := `{"success": false, "message": "Invalid event ID."}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// No need to check service expectations since it shouldn't be called
}

func TestDeleteEvent_ServiceError(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.RTOBLL)

	// Mock the service method to return an error
	mockService.On("TransformVacationToRemote", 2).Return(errors.New("event not found"))

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock("none", mockService, QuarterStart, QuarterEnd)
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger

	// Create a DELETE request with a valid event ID but service returns error
	req := httptest.NewRequest(http.MethodDelete, "/events/delete/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("2")

	// Call the handler
	if assert.NoError(t, ctlr.DeleteEvent(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		expectedResponse := `{"success": false, "message": "event not found"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	// Verify that the expectations were met
	mockService.AssertExpectations(t)
}
