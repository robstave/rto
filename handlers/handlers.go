package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// AddEvent handles the addition of new events via form submission
func AddEvent(c echo.Context) error {
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
		log.Printf("Error parsing date: %v", err)
		return c.String(http.StatusBadRequest, "Invalid date format")
	}

	// Initialize Event struct
	newEvent := Event{
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

	// Append the new event to allEvents
	allEvents = append(allEvents, newEvent)

	// Save to events.json
	/*
		file, err := os.Create("data/events.json")
		if err != nil {
			log.Printf("Error creating events.json: %v", err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(allEvents); err != nil {
			log.Printf("Error encoding events to JSON: %v", err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	*/

	// Redirect back to the calendar
	return c.Redirect(http.StatusSeeOther, "/")
}

// ShowAddEventForm renders the Add Event form
func ShowAddEventForm(c echo.Context) error {
	return c.Render(http.StatusOK, "add_event.html", nil)
}

// EventsList handles displaying the list of events
func EventsList(c echo.Context) error {
	// Pass allEvents to the template
	data := map[string]interface{}{
		"Events": allEvents,
	}

	return c.Render(http.StatusOK, "events.html", data)
}
