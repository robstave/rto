package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/internal/domain/types"
)

// EventsList handles displaying the list of events
func (ctlr *RTOController) EventsList(c echo.Context) error {
	// Pass allEvents to the template
	data := map[string]interface{}{
		"Events": ctlr.service.GetAllEvents(),
	}

	return c.Render(http.StatusOK, "events.html", data)
}

// ShowAddEventForm renders the Add Event form
func (ctlr *RTOController) ShowAddEventForm(c echo.Context) error {
	return c.Render(http.StatusOK, "add_event.html", nil)
}

// AddEvent handles the addition of new events via form submission
func (ctlr *RTOController) AddEvent(c echo.Context) error {
	dateStr := c.FormValue("date")   // Expected format: YYYY-MM-DD
	eventType := c.FormValue("type") // "holiday", "vacation", "attendance"
	description := c.FormValue("description")
	isInOfficeStr := c.FormValue("isInOffice") // "true" or "false"

	if dateStr == "" || eventType == "" {
		return c.String(http.StatusBadRequest, "Date and Event Type are required")
	}

	// Parse the date
	eventDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctlr.logger.Error("Error parsing date", "fn", "AddEvent", "date", err)
		return c.String(http.StatusBadRequest, "Invalid date format")
	}

	// Initialize Event struct
	newEvent := types.Event{
		Date:        eventDate,
		Description: description,
		Type:        eventType,
	}

	// Handle Attendance Type
	if eventType == "attendance" {
		if isInOfficeStr == "true" {
			newEvent.IsInOffice = true
		} else {
			newEvent.IsInOffice = false
		}
	}

	// Call domain service to add event
	err = ctlr.service.AddEvent(newEvent)
	if err != nil {
		ctlr.logger.Error("Error adding event", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add event.")
	}

	// Redirect back to the calendar
	return c.Redirect(http.StatusSeeOther, "/")
}

func (ctlr *RTOController) AddDefaultDays(c echo.Context) error {
	err := ctlr.service.AddDefaultDays()
	if err != nil {
		ctlr.logger.Error("Error adding default days", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add default attendance events.")
	}
	return c.Redirect(http.StatusSeeOther, "/prefs")
}

type RawEvent struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// BulkAddEventsJSONRequest represents the expected JSON payload for bulk adding events
type BulkAddEventsJSONRequest struct {
	Events []RawEvent `json:"events"`
}

// BulkAddEventsJSON handles the bulk addition of vacation events via JSON
func (ctlr *RTOController) BulkAddEventsJSON(c echo.Context) error {
	var rawEventsReq BulkAddEventsJSONRequest

	// Parse and decode the JSON request body into rawEvents
	if err := c.Bind(&rawEventsReq); err != nil {
		ctlr.logger.Error("Error binding bulk add JSON", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid JSON payload.",
		})
	}

	var domainEvents []types.Event

	// Iterate over each raw event to validate and transform
	for i, rawEvent := range rawEventsReq.Events {
		// Basic validation
		if strings.TrimSpace(rawEvent.Date) == "" ||
			strings.TrimSpace(rawEvent.Description) == "" ||
			strings.TrimSpace(rawEvent.Type) == "" {
			ctlr.logger.Error("Incomplete event data in bulk add", "event_index", i, "event", rawEvent)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("Event at index %d is missing required fields.", i),
			})
		}

		// Ensure the event type is 'vacation'
		if strings.ToLower(rawEvent.Type) != "vacation" {
			ctlr.logger.Error("Invalid event type in bulk add", "event_index", i, "event", rawEvent)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("Only events with type 'vacation' can be bulk added (error at index %d).", i),
			})
		}

		// Validate and parse date (YYYY-MM-DD)
		parsedDate, err := time.Parse("2006-01-02", rawEvent.Date)
		if err != nil {
			ctlr.logger.Error("Invalid date format in bulk add", "event_index", i, "event", rawEvent, "error", err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("Invalid date format for event on %s. Expected YYYY-MM-DD.", rawEvent.Date),
			})
		}

		// Create domain Event
		domainEvent := types.Event{
			Date:        parsedDate, // Assuming types.Date is defined in the domain layer
			Description: rawEvent.Description,
			Type:        strings.ToLower(rawEvent.Type),
		}

		domainEvents = append(domainEvents, domainEvent)
	}

	// Delegate the processing to the service layer
	response, err := ctlr.service.BulkAddEvents(domainEvents)
	if err != nil {
		ctlr.logger.Error("Error in BulkAddEvents service method", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "An error occurred while processing bulk add events.",
		})
	}

	// Return the service layer's response
	return c.JSON(http.StatusOK, response)
}
