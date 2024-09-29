package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"os"
	"time"
)

// Event represents a calendar event
type Event struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	IsInOffice  bool      `json:"isInOffice"`
	Type        string    `json:"type"` // "attendance", "holiday", "vacation"
}

// InitializeEvents loads holidays and attendance events without duplicates
func InitializeEvents(holidaysPath string, eventsPath string) {
	holidays, err := LoadHolidays(holidaysPath)
	if err != nil {
		logger.Error("Error loading holidays", "error", err)
		allEvents = []Event{}
		return
	}

	// Initialize a map to track unique events
	eventMap := make(map[string]Event)

	// Add holidays to the map
	for _, h := range holidays {
		key := h.Date.Format("2006-01-02") + "_" + h.Type
		eventMap[key] = h
	}

	// Load attendance events from events.json
	attendanceEvents, err := LoadAttendanceEvents(eventsPath)
	if err != nil {
		logger.Error("Error loading attendance events:",
			"error", err)
		attendanceEvents = []Event{}
	}

	// Add attendance events to the map, avoiding duplicates
	for _, a := range attendanceEvents {
		key := a.Date.Format("2006-01-02") + "_" + a.Type
		if _, exists := eventMap[key]; !exists {
			eventMap[key] = a
		} else {
			logger.Info(
				"Duplicate event found. Skipping.",
				"date", a.Date.Format("2006-01-02"),
				"type", a.Type,
			)
		}
	}

	// Convert the map back to a slice
	allEvents = []Event{}
	for _, e := range eventMap {
		allEvents = append(allEvents, e)
	}

	logger.Info("Events Loaded", "len", len(allEvents))
}

// LoadAttendanceEvents loads attendance events from a JSON file
func LoadAttendanceEvents(filePath string) ([]Event, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var rawEvents2 []struct {
		Date        string `json:"date"`
		Description string `json:"description"`
		Type        string `json:"type"`
		IsInOffice  bool   `json:"isInOffice"`
	}

	if err := json.Unmarshal(byteValue, &rawEvents2); err != nil {
		return nil, err
	}
	// Ensure date parsing
	var events []Event
	for _, e := range rawEvents2 {

		parsedDate, err := parseDate(e.Date)
		if err != nil {
			logger.Error(
				"Invalid date format in events.json",
				"event", e.Description,
				"error", err)
			continue
		}

		events = append(events, Event{
			Date:        parsedDate,
			Description: e.Description,
			IsInOffice:  e.IsInOffice, // Holidays and vacations override attendance
			Type:        e.Type,
		})
	}

	return events, nil
}

// RawHoliday represents the structure of each holiday entry in the JSON file.
type RawHoliday struct {
	Date string `json:"date"`
	Name string `json:"name"`
	Type string `json:"type"` // e.g., "holiday", "vacation"
}

// LoadHolidays loads holidays and vacations from a JSON file
func LoadHolidays(filePath string) ([]Event, error) {

	logger.Info("loading Holidays")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var rawEvents []RawHoliday

	if err := json.Unmarshal(byteValue, &rawEvents); err != nil {
		return nil, err
	}

	events, processingErrors := processRawHolidays(rawEvents)

	if len(processingErrors) > 0 {
		logger.Error("Encountered errors while processing holidays.", "len", len(processingErrors))
	}

	return events, nil
}

// processRawHolidays converts raw holiday data into Event structs.
// It returns a slice of Events and a slice of errors encountered during processing.
func processRawHolidays(rawEvents []RawHoliday) ([]Event, []error) {
	events := []Event{}
	errorsList := []error{}

	for _, re := range rawEvents {
		parsedDate, err := parseDate(re.Date)
		if err != nil {
			logger.Error("Invalid date format in holidays.json", "error", err)
			errorsList = append(errorsList, err)
			continue
		}
		events = append(events, Event{
			Date:        parsedDate,
			Description: re.Name,
			IsInOffice:  false, // Holidays and vacations override attendance
			Type:        re.Type,
		})
	}

	return events, errorsList
}

// parseDate tries to parse a date string using multiple layouts.
// It returns the parsed time.Time or an error if none of the layouts match.
func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",          // "YYYY-MM-DD"
		time.RFC3339,          // "YYYY-MM-DDTHH:MM:SSZ"
		"2006-01-02T15:04:05", // "YYYY-MM-DDTHH:MM:SS"
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("invalid date format")
}
