package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"
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

// AddEvent adds a new event or updates an existing one based on the event type and date
func (s *Service) AddEvent(event types.Event) error {

	event.Date = utils.NormalizeDate(event.Date)

	if event.Type == "vacation" {
		// Check if a vacation event already exists on the given date
		existingEvent, err := s.eventRepo.GetEventByDateAndType(event.Date, "vacation")
		if err != nil && err != gorm.ErrRecordNotFound {
			s.logger.Error("Error fetching existing vacation event", "error", err)
			return err
		}

		if existingEvent.ID != 0 {
			// Vacation event exists; update the description
			existingEvent.Description = event.Description
			err = s.eventRepo.UpdateEvent(existingEvent)
			if err != nil {
				s.logger.Error("Error updating vacation event", "error", err)
				return err
			}
			s.logger.Info("Vacation event updated", "date", event.Date)
			return nil
		}
	} else if event.Type == "attendance" {
		s.logger.Info("ADding Attendence", "date", event.Date, "type", event.Type)

		// Check if an attendance event already exists on the given date
		existingEvent, err := s.eventRepo.GetEventByDateAndType(event.Date, "attendance")
		if err != nil && err != gorm.ErrRecordNotFound {
			s.logger.Error("Error fetching existing attendance event", "error", err)
			return err
		}

		s.logger.Info("doing add", "existingEvent", existingEvent.ID)

		if existingEvent.ID != 0 {
			// Attendance event exists; do not add a duplicate
			s.logger.Info("Attendance event already exists", "date", event.Date.Format("2006-01-02"))
			return nil // Alternatively, return a custom error if you want to notify the controller
		}
	}

	s.logger.Info("doing add", "date", event.Date, "type", event.Type)

	// No existing event; proceed to add the new event
	err := s.eventRepo.AddEvent(event)
	if err != nil {
		s.logger.Error("Error adding event", "error", err)
		return err
	}
	s.logger.Info("Event added", "date", event.Date.Format("2006-01-02"), "type", event.Type)
	return nil
}

// ClearEventsForDate clears all events for a specific date
func (s *Service) ClearEventsForDate(date time.Time) error {
	events, err := s.eventRepo.GetEventsByDate(date)
	s.logger.Info("0000-----ClearEventsForDate------", "date", date, "len", len(events))

	if err != nil {
		s.logger.Error("Error fetching events for date", "date", date, "error", err)
		return err
	}

	for _, event := range events {
		s.logger.Info("000bbbb0-----deletin------", "len", int(event.ID))

		err := s.eventRepo.DeleteEvent(int(event.ID))
		if err != nil {
			s.logger.Error("Error deleting event", "eventID", event.ID, "error", err)
			return err
		}
	}

	s.logger.Info("All events cleared for date", "date", date.Format("2006-01-02"))
	return nil
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
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.UTC)

	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.UTC)

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
	_, err := s.GetEventByID(eventID)
	if err != nil {
		return err
	}

	////if event.Type != "vacation" {
	//return errors.New("only vacation events can be deleted")
	//}

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

// GetEventsByDate retrieves all events for a specific date
func (s *Service) GetEventsByDate(date time.Time) ([]types.Event, error) {
	events, err := s.eventRepo.GetEventsByDate(date)
	if err != nil {
		s.logger.Error("Error fetching events by date", "date", date, "error", err)
		return nil, err
	}
	return events, nil
}
