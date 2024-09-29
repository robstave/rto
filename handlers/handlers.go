package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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
		logger.Error("Error parsing date", "fn", "AddEvent", "date", err)
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

	// Redirect back to the calendar
	return c.Redirect(http.StatusSeeOther, "/")
}

// ToggleAttendanceRequest represents the JSON payload for toggling attendance
type ToggleAttendanceRequest struct {
	Date string `json:"date"` // Expected format: YYYY-MM-DD
}

// ToggleAttendanceResponse represents the JSON response after toggling
type ToggleAttendanceResponse struct {
	Success       bool    `json:"success"`
	NewStatus     string  `json:"newStatus,omitempty"` // "in" or "remote"
	Message       string  `json:"message,omitempty"`
	InOfficeCount int     `json:"inOfficeCount,omitempty"`
	TotalDays     int     `json:"totalDays,omitempty"`
	Average       float64 `json:"average,omitempty"`
	AverageDays   float64 `json:"averageDays,omitempty"`
	TargetDays    float64 `json:"targetDays,omitempty"` // New field for target value
}

// ToggleAttendance handles toggling attendance status for a given date
func ToggleAttendance(c echo.Context) error {

	req := new(ToggleAttendanceRequest)
	if err := c.Bind(req); err != nil {
		logger.Error("Error binding request: %v", "error", err)
		return c.JSON(http.StatusBadRequest, ToggleAttendanceResponse{
			Success: false,
			Message: "Invalid request payload.",
		})
	}
	logger.Info("togglr yo", "date", req.Date)
	// Parse the date
	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		logger.Error("Error parsing date", "fn", "Toggle Attendance", "error", err)
		return c.JSON(http.StatusBadRequest, ToggleAttendanceResponse{
			Success: false,
			Message: "Invalid date format. Expected YYYY-MM-DD.",
		})
	}

	// Lock events for writing
	eventsLock.Lock()
	defer eventsLock.Unlock()
	logger.Info("searching events", "size", len(allEvents))
	// Find the attendance event on the given date
	found := false
	for i, event := range allEvents {
		if sameDay(event.Date, eventDate) && event.Type == "attendance" {
			// Toggle the IsInOffice flag
			allEvents[i].IsInOffice = !event.IsInOffice
			logger.Info("found and toggled", "value", allEvents[i].IsInOffice)
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

	// Save to events.json
	eventsFilePath := "data/events.json" // Ensure this path is correct
	if err := SaveEvents(eventsFilePath); err != nil {
		logger.Error("Error saving events", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to save event.")
	}

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

	// Calculate updated counts and averages
	currentYear := time.Now().Year()
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.Local)

	inOfficeCount, totalDays := calculateInOfficeAverage(allEvents, startDate, endDate)

	average := 0.0
	averageDays := 0.0
	if totalDays > 0 {
		average = (float64(inOfficeCount) / float64(totalDays)) * 100
		averageDays = (float64(inOfficeCount) / float64(totalDays)) * 7 // Average days/week
	}

	// Fetch target days from preferences
	preferencesLock.RLock()
	targetDaysStr := preferences.TargetDays // Assuming TargetDays is added to Preferences
	preferencesLock.RUnlock()

	targetDays, err := strconv.ParseFloat(targetDaysStr, 64)
	if err != nil {
		// Fallback to default target if parsing fails
		targetDays = 2.5
	}

	return c.JSON(http.StatusOK, ToggleAttendanceResponse{
		Success:       true,
		NewStatus:     newStatus,
		InOfficeCount: inOfficeCount,
		TotalDays:     totalDays,
		Average:       average,
		AverageDays:   averageDays,
		TargetDays:    targetDays,
	})
}

// ShowPrefs renders the preferences page with current default in-office days and target
func ShowPrefs(c echo.Context) error {
	preferencesLock.RLock()
	defer preferencesLock.RUnlock()

	data := map[string]interface{}{
		"Preferences": preferences,
	}

	return c.Render(http.StatusOK, "prefs.html", data)
}

// UpdatePreferences handles updating user preferences
func UpdatePreferences(c echo.Context) error {
	newDefaultDays := c.FormValue("defaultDays")
	newTargetDays := c.FormValue("targetDays")

	if newDefaultDays == "" || newTargetDays == "" {
		return c.String(http.StatusBadRequest, "Default Days and Target Days are required.")
	}

	// Update preferences
	preferencesLock.Lock()
	preferences.DefaultDays = newDefaultDays
	preferences.TargetDays = newTargetDays
	preferencesLock.Unlock()

	// Save preferences to JSON file
	preferencesFilePath := "data/preferences.json"
	if err := SavePreferences(preferencesFilePath); err != nil {
		logger.Error("Error saving preferences", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to save preferences.")
	}

	return c.Redirect(http.StatusSeeOther, "/prefs")
}

// SavePreferences saves the current preferences to the specified JSON file.
func SavePreferences(filePath string) error {
	preferencesLock.RLock()
	defer preferencesLock.RUnlock()

	data, err := json.MarshalIndent(preferences, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
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
	logger.Info("AddDefaultDays triggered")

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

		// **New Condition:** Skip Saturdays and Sundays
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue // Skip to the next day without creating an event
		}

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

	logger.Info("AddDefaultDays: added events.", "count", addedCount)

	// Optionally, you can save `allEvents` to `events.json` here if persistence is desired

	// Save the updated events to events.json
	eventsFilePath := "data/events.json" // Ensure this path is correct
	if err := SaveEvents(eventsFilePath); err != nil {
		logger.Error("Error saving events", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to save default attendance events.")
	}

	// Redirect back to preferences with a success message
	// To display messages, you can use query parameters or session-based flash messages
	// For simplicity, we'll redirect without messages
	return c.Redirect(http.StatusSeeOther, "/prefs")
}

// SaveEvents saves the current list of events to the specified JSON file.
func SaveEvents(filePath string) error {

	//eventsLock.RLock()
	//defer eventsLock.RUnlock()

	data, err := json.MarshalIndent(allEvents, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
