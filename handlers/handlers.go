package handlers

import (
	"log"
	"net/http"
	"strings"
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

// ToggleAttendanceRequest represents the JSON payload for toggling attendance
type ToggleAttendanceRequest struct {
	Date string `json:"date"` // Expected format: YYYY-MM-DD
}

// ToggleAttendanceResponse represents the JSON response after toggling
type ToggleAttendanceResponse struct {
	Success   bool   `json:"success"`
	NewStatus string `json:"newStatus,omitempty"` // "in" or "remote"
	Message   string `json:"message,omitempty"`
}

// ToggleAttendance handles toggling attendance status for a given date
func ToggleAttendance(c echo.Context) error {

	req := new(ToggleAttendanceRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, ToggleAttendanceResponse{
			Success: false,
			Message: "Invalid request payload.",
		})
	}
	log.Println("togglr yo", req.Date)
	// Parse the date
	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return c.JSON(http.StatusBadRequest, ToggleAttendanceResponse{
			Success: false,
			Message: "Invalid date format. Expected YYYY-MM-DD.",
		})
	}

	// Lock events for writing
	eventsLock.Lock()
	defer eventsLock.Unlock()
	log.Println("searching events len:", len(allEvents))
	// Find the attendance event on the given date
	found := false
	for i, event := range allEvents {
		log.Println("x", event.Type, event.Date, eventDate)
		if sameDay(event.Date, eventDate) && event.Type == "attendance" {
			// Toggle the IsInOffice flag
			allEvents[i].IsInOffice = !event.IsInOffice
			log.Println("found and toggled", allEvents[i].IsInOffice)
			found = true
			break
		}
	}

	if !found {
		return c.JSON(http.StatusNotFound, ToggleAttendanceResponse{
			Success: false,
			Message: "Attendance event not found on the specified date.",
		})
	}

	// Optionally, save the updated events to events.json for persistence
	// Uncomment the following lines if you wish to persist changes

	/*
		err = SaveEvents("data/events.json")
		if err != nil {
			log.Printf("Error saving events: %v", err)
			return c.JSON(http.StatusInternalServerError, ToggleAttendanceResponse{
				Success: false,
				Message: "Failed to save updated events.",
			})
		}
	*/

	// Determine the new status
	newStatus := "remote"
	for _, event := range allEvents {
		if sameDay(event.Date, eventDate) && event.Type == "attendance" {
			if event.IsInOffice {
				newStatus = "in"
			}
			break
		}
	}

	log.Println("done")
	return c.JSON(http.StatusOK, ToggleAttendanceResponse{
		Success:   true,
		NewStatus: newStatus,
	})
}

// ShowPrefs renders the preferences page with current default in-office days
func ShowPrefs(c echo.Context) error {
	preferencesLock.RLock()
	defer preferencesLock.RUnlock()

	data := map[string]interface{}{
		"Preferences": preferences,
	}

	return c.Render(http.StatusOK, "prefs.html", data)
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

// AddDefaultDays handles the addition of default attendance events
func AddDefaultDays(c echo.Context) error {
	log.Println("AddDefaultDays triggered")

	// Define the date range: October 1 to December 31 of the current year
	currentYear := time.Now().Year()
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.Local)

	// Retrieve default in-office days from preferences
	preferencesLock.RLock()
	defaultDays := strings.Split(preferences.DefaultDays, ",")
	preferencesLock.RUnlock()

	// Create a map for faster lookup of default in-office days
	defaultDaysMap := make(map[string]bool)
	for _, day := range defaultDays {
		day = strings.TrimSpace(strings.ToLower(day))
		defaultDaysMap[day] = true
	}

	// Lock events for writing
	eventsLock.Lock()
	defer eventsLock.Unlock()

	// Counter for added events
	addedCount := 0

	// Iterate through each day in the date range
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		//weekday := strings.ToLower(d.Weekday().String()) // e.g., "monday"

		// Map Go's Weekday to user's day abbreviations
		var dayAbbrev string
		switch d.Weekday() {
		case time.Monday:
			dayAbbrev = "M"
		case time.Tuesday:
			dayAbbrev = "T"
		case time.Wednesday:
			dayAbbrev = "W"
		case time.Thursday:
			dayAbbrev = "Th"
		case time.Friday:
			dayAbbrev = "F"
		case time.Saturday:
			dayAbbrev = "Sat"
		case time.Sunday:
			dayAbbrev = "Sun"
		}

		dayAbbrevLower := strings.ToLower(dayAbbrev)

		// Determine if it's a default in-office day
		isInOffice, isDefault := defaultDaysMap[dayAbbrevLower]
		if !isDefault {
			// Non-default days are considered remote
			isInOffice = false
		}

		// Check if an event already exists on this day
		eventExists := false
		for _, event := range allEvents {
			if sameDay(event.Date, d) {
				eventExists = true
				break
			}
		}

		if !eventExists {
			// Create a new attendance event
			newEvent := Event{
				Date:        d,
				Description: "",
				IsInOffice:  isInOffice,
				Type:        "attendance",
			}
			allEvents = append(allEvents, newEvent)
			addedCount++
		}
	}

	log.Printf("AddDefaultDays: Added %d new attendance events.", addedCount)

	// Optionally, you can save `allEvents` to `events.json` here if persistence is desired

	// Redirect back to preferences with a success message
	// To display messages, you can use query parameters or session-based flash messages
	// For simplicity, we'll redirect without messages
	return c.Redirect(http.StatusSeeOther, "/prefs")
}

// sameDay checks if two dates are on the same calendar day
func sameDay(a, b time.Time) bool {
	yearA, monthA, dayA := a.Date()
	yearB, monthB, dayB := b.Date()
	return yearA == yearB && monthA == monthB && dayA == dayB
}
