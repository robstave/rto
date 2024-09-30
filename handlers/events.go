package handlers

import (
	"encoding/json"
	"io/ioutil"

	"os"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"
)

// InitializeEvents loads holidays and attendance events without duplicates
func InitializeEvents(holidaysPath string, eventsPath string) {
	holidays, err := LoadHolidays(holidaysPath)
	if err != nil {
		logger.Error("Error loading holidays", "error", err)
		allEvents = []types.Event{}
		return
	}

	// Initialize a map to track unique events
	eventMap := make(map[string]types.Event)

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
		attendanceEvents = []types.Event{}
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
	allEvents = []types.Event{}
	for _, e := range eventMap {
		allEvents = append(allEvents, e)
	}

	logger.Info("Events Loaded", "len", len(allEvents))
}

// LoadAttendanceEvents loads attendance events from a JSON file
func LoadAttendanceEvents(filePath string) ([]types.Event, error) {
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
	var events []types.Event
	for _, e := range rawEvents2 {

		parsedDate, err := utils.ParseDate(e.Date)
		if err != nil {
			logger.Error(
				"Invalid date format in events.json",
				"event", e.Description,
				"error", err)
			continue
		}

		events = append(events, types.Event{
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
func LoadHolidays(filePath string) ([]types.Event, error) {

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
func processRawHolidays(rawEvents []RawHoliday) ([]types.Event, []error) {
	events := []types.Event{}
	errorsList := []error{}

	for _, re := range rawEvents {
		parsedDate, err := utils.ParseDate(re.Date)
		if err != nil {
			logger.Error("Invalid date format in holidays.json", "error", err)
			errorsList = append(errorsList, err)
			continue
		}
		events = append(events, types.Event{
			Date:        parsedDate,
			Description: re.Name,
			IsInOffice:  false, // Holidays and vacations override attendance
			Type:        re.Type,
		})
	}

	return events, errorsList
}
