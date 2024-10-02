package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/internal/domain/mocks"
	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestEventsList(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a mock RTOBLL
	mockService := new(mocks.MockRTOBLL)

	// Define mock events
	mockEvents := []types.Event{
		{
			ID:          1,
			Date:        time.Now(),
			Description: "Company Holiday",
			IsInOffice:  false,
			Type:        "holiday",
		},
		{
			ID:          2,
			Date:        time.Now(),
			Description: "Vacation",
			IsInOffice:  false,
			Type:        "vacation",
		},
		{
			ID:          3,
			Date:        time.Now(),
			Description: "",
			IsInOffice:  true,
			Type:        "attendance",
		},
	}

	// Setup expectations
	mockService.On("GetAllEvents").Return(mockEvents)

	// Initialize the controller with the mock service
	ctlr := NewRTOControllerWithMock(mockService)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set renderer (mock or real)
	// For simplicity, we'll use a minimal renderer
	ctlr.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Assign a simple logger
	e.Renderer = &mockRenderer{}

	// Call the handler
	if assert.NoError(t, ctlr.EventsList(c)) {
		// Assertions on the response
		assert.Equal(t, http.StatusOK, rec.Code)
		// Further assertions can be made on the rendered content if using a real renderer
	}

	// Assert that the expectations were met
	mockService.AssertExpectations(t)
}

// mockRenderer is a minimal implementation of echo.Renderer for testing purposes
type mockRenderer struct{}

func (m *mockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// For testing, we can simply write a success message or serialize the data
	// Here, we'll serialize the data as JSON for simplicity
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	return err
}
