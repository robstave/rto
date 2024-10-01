package domain

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"os"
	"strings"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"
)

func (s *Service) GetAllEvents() []types.Event {
	events, err := s.eventRepo.GetAllEvents()
	if err != nil {
		s.logger.Error("Error getting events", "error", err)
		return []types.Event{}
	}
	return events
}
func (s *Service) GetPrefs() types.Preferences {
	prefs, err := s.preferenceRepo.GetPreferences()
	if err != nil {
		s.logger.Error("Error getting preferences", "error", err)
		// Return default preferences or handle error as needed
	}
	return prefs
}

func (s *Service) AddEvent(event types.Event) error {
	err := s.eventRepo.AddEvent(event)
	if err != nil {
		s.logger.Error("Error adding event", "error", err)
	}
	return err
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

// InitializeEvents loads holidays and attendance events without duplicates
func initializeEvents(s *Service, holidaysPath string, eventsPath string) {
	holidays, err := LoadHolidays(s, holidaysPath)
	if err != nil {
		s.logger.Error("Error loading holidays", "error", err)
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
	attendanceEvents, err := loadAttendanceEvents(s, eventsPath)
	if err != nil {
		s.logger.Error("Error loading attendance events:",
			"error", err)
		attendanceEvents = []types.Event{}
	}

	// Add attendance events to the map, avoiding duplicates
	for _, a := range attendanceEvents {
		key := a.Date.Format("2006-01-02") + "_" + a.Type
		if _, exists := eventMap[key]; !exists {
			eventMap[key] = a
		} else {
			s.logger.Info(
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

	s.logger.Info("Events Loaded", "len", len(allEvents))
}

// LoadAttendanceEvents loads attendance events from a JSON file
func loadAttendanceEvents(s *Service, filePath string) ([]types.Event, error) {
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
			s.logger.Error(
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
func LoadHolidays(s *Service, filePath string) ([]types.Event, error) {

	s.logger.Info("loading Holidays")
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

	events, processingErrors := processRawHolidays(s, rawEvents)

	if len(processingErrors) > 0 {
		s.logger.Error("Encountered errors while processing holidays.", "len", len(processingErrors))
	}

	return events, nil
}

// processRawHolidays converts raw holiday data into Event structs.
// It returns a slice of Events and a slice of errors encountered during processing.
func processRawHolidays(s *Service, rawEvents []RawHoliday) ([]types.Event, []error) {
	events := []types.Event{}
	errorsList := []error{}

	for _, re := range rawEvents {
		parsedDate, err := utils.ParseDate(re.Date)
		if err != nil {
			s.logger.Error("Invalid date format in holidays.json", "error", err)
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

func (s *Service) AddDefaultDays() error {
	s.logger.Info("AddDefaultDays triggered")

	// Define the date range: October 1 to December 31 of the current year
	currentYear := time.Now().Year()
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.Local)

	// Retrieve default in-office days from preferences
	//s.preferencesLock.RLock()
	defaultDays := strings.Split(s.preferences.DefaultDays, ",")
	//s.preferencesLock.RUnlock()

	// Create a map for faster lookup of default in-office days
	defaultDaysMap := make(map[string]bool)
	for _, day := range defaultDays {
		day = strings.TrimSpace(strings.ToLower(day))
		defaultDaysMap[day] = true
	}

	// Lock events for writing
	//s.eventsLock.Lock()
	//defer s.eventsLock.Unlock()

	// Counter for added events
	addedCount := 0

	// Iterate through each day in the date range
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
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

		// Skip weekends
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
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
			if utils.SameDay(event.Date, d) {
				eventExists = true
				break
			}
		}

		if !eventExists {
			// Create a new attendance event
			newEvent := types.Event{
				Date:        d,
				Description: "",
				IsInOffice:  isInOffice,
				Type:        "attendance",
			}
			allEvents = append(allEvents, newEvent)
			addedCount++
		}
	}

	s.logger.Info("AddDefaultDays: added events.", "count", addedCount)

	// Save the updated events to events.json
	eventsFilePath := "data/events.json"
	if err := SaveEvents(eventsFilePath); err != nil {
		s.logger.Error("Error saving events", "error", err)
		return err
	}

	return nil
}
