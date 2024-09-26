package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AddEvent handles the addition of new events
func AddEvent(c echo.Context) error {
	date := c.FormValue("date") // Expected format: YYYY-MM-DD
	event := c.FormValue("event")

	if date == "" || event == "" {
		return c.String(http.StatusBadRequest, "Date and Event are required")
	}

	// Here, you'd typically save the event to a database
	// For demonstration, we'll append it to the existing events map
	// Note: Since getSampleEvents returns a new map each time, you'd need a persistent storage

	// Example (In-Memory - Not Persistent)
	// This requires modifying getSampleEvents to use a global variable or a proper storage mechanism

	return c.Redirect(http.StatusSeeOther, "/")
}
