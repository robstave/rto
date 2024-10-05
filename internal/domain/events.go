package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"gorm.io/gorm"
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

func (s *Service) AddDefaultDays() error {
	s.logger.Info("AddDefaultDays triggered")

	// Get current preferences
	prefs, err := s.preferenceRepo.GetPreferences()
	if err != nil {
		s.logger.Error("Failed to get preferences", "error", err)
		return err
	}

	// Parse default days
	defaultDays := strings.Split(prefs.DefaultDays, ",")
	defaultDaysMap := make(map[string]bool)
	for _, day := range defaultDays {
		day = strings.TrimSpace(strings.ToLower(day))
		defaultDaysMap[day] = true
	}

	// Define the date range
	currentYear := time.Now().Year()
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.Local)

	// Retrieve existing events
	existingEvents, err := s.eventRepo.GetAllEvents()
	if err != nil {
		s.logger.Error("Failed to retrieve events", "error", err)
		return err
	}

	// Create a map of existing event dates
	existingEventDates := make(map[string]bool)
	for _, event := range existingEvents {
		dateStr := event.Date.Format("2006-01-02")
		existingEventDates[dateStr] = true
	}

	// Add default attendance events
	addedCount := 0
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
		isInOffice := defaultDaysMap[dayAbbrevLower]

		dateStr := d.Format("2006-01-02")
		if !existingEventDates[dateStr] {
			// Create a new attendance event
			newEvent := types.Event{
				Date:        d,
				Description: "",
				IsInOffice:  isInOffice,
				Type:        "attendance",
			}
			err := s.eventRepo.AddEvent(newEvent)
			if err != nil {
				s.logger.Error("Failed to add event", "date", dateStr, "error", err)
				continue
			}
			addedCount++
		}
	}

	s.logger.Info("AddDefaultDays completed", "events_added", addedCount)
	return nil
}

// GetEventByID retrieves a single event by its ID
func (s *Service) GetEventByID(eventID int) (types.Event, error) {
	event, err := s.eventRepo.GetEventByID(eventID)
	if err != nil {
		s.logger.Error("Error fetching event by ID", "error", err)
		return types.Event{}, err
	}
	return event, nil
}

// DeleteEvent deletes an event by its ID
func (s *Service) DeleteEvent(eventID int) error {
	// First, retrieve the event to ensure it exists and is deletable
	event, err := s.GetEventByID(eventID)
	if err != nil {
		return err
	}

	if event.Type != "vacation" {
		return errors.New("only vacation events can be deleted")
	}

	// Proceed to delete the event
	err = s.eventRepo.DeleteEvent(eventID)
	if err != nil {
		s.logger.Error("Error deleting event", "error", err)
		return err
	}

	return nil
}

// UpdateEvent updates an existing event in the database
func (s *Service) UpdateEvent(event types.Event) error {
	if event.ID == 0 {
		return errors.New("event ID is required for update")
	}
	err := s.eventRepo.UpdateEvent(event)
	if err != nil {
		s.logger.Error("Failed to update event", "eventID", event.ID, "error", err)
		return err
	}
	return nil
}

// GetEventByDateAndType retrieves an event by date and type
func (s *Service) GetEventByDateAndType(date time.Time, eventType string) (*types.Event, error) {
	event, err := s.eventRepo.GetEventByDateAndType(date, eventType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}
	return &event, nil
}
