package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Event represents a calendar event
type Event struct {
	Date time.Time `json:"date"`
	//Date string `json:"date"`

	Description string `json:"description"`
	IsInOffice  bool   `json:"isInOffice"`
	Type        string `json:"type"` // "attendance", "holiday", "vacation"
}

// InitializeEvents loads holidays and attendance events without duplicates
func InitializeEvents(holidaysPath string, eventsPath string) {
	holidays, err := LoadHolidays(holidaysPath)
	if err != nil {
		log.Printf("Error loading holidays: %v", err)
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
		log.Printf("Error loading attendance events: %v", err)
		attendanceEvents = []Event{}
	}

	// Add attendance events to the map, avoiding duplicates
	for _, a := range attendanceEvents {
		key := a.Date.Format("2006-01-02") + "_" + a.Type
		if _, exists := eventMap[key]; !exists {
			eventMap[key] = a
		} else {
			log.Printf("Duplicate event found on %s with type %s. Skipping.", a.Date.Format("2006-01-02"), a.Type)
		}
	}

	// Convert the map back to a slice
	allEvents = []Event{}
	for _, e := range eventMap {
		allEvents = append(allEvents, e)
	}

	log.Printf("Events Loaded: len= %v", len(allEvents))
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

	//var rawEvents2 []Event
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
	for i, e := range rawEvents2 {
		log.Println("1")
		//parsedDate, err := time.Parse("2006-01-02", e.Date.Format("2006-01-02"))
		//parsedDate, err := time.Parse("2006-01-02", e.Date)
		parsedDate, err := parseDate(e.Date)
		if err != nil {
			log.Printf("Invalid date format in events.json %v for event %v: %v", i, e.Description, err)
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

// LoadHolidays loads holidays and vacations from a JSON file
func LoadHolidays(filePath string) ([]Event, error) {

	log.Println("loading Holidays")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var rawEvents []struct {
		Date string `json:"date"`
		Name string `json:"name"`
		Type string `json:"type"`
	}

	if err := json.Unmarshal(byteValue, &rawEvents); err != nil {
		return nil, err
	}

	events := []Event{}
	for _, re := range rawEvents {
		parsedDate, err := parseDate(re.Date)
		if err != nil {
			log.Printf("Invalid date format in holidays.json: %v", err)
			continue
		}
		events = append(events, Event{
			Date:        parsedDate,
			Description: re.Name,
			IsInOffice:  false, // Holidays and vacations override attendance
			Type:        re.Type,
		})
	}

	return events, nil
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
