package repository

import (
	"time"

	"github.com/robstave/rto/internal/domain/types"
)

// Existing EventRepository interface methods...

type EventRepository interface {
	GetAllEvents() ([]types.Event, error)
	AddEvent(event types.Event) error
	UpdateEvent(event types.Event) error
	DeleteEvent(eventID int) error
	GetEventByDate(date time.Time) (types.Event, error)
	GetEventByID(eventID int) (types.Event, error)
	GetEventsByType(eventType string) ([]types.Event, error)

	// New method to fetch event by date and type
	GetEventByDateAndType(date time.Time, eventType string) (types.Event, error)
}

type PreferenceRepository interface {
	GetPreferences() (types.Preferences, error)
	UpdatePreferences(prefs types.Preferences) error
	// Add other methods as needed
}
