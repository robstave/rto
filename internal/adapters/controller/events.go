package controller

import (
	"net/http"
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
