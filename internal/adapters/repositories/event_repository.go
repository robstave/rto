package repository

import (
	"time"

	"github.com/robstave/rto/internal/domain/types"
)

type EventRepository interface {
	GetAllEvents() ([]types.Event, error)
	AddEvent(event types.Event) error
	UpdateEvent(event types.Event) error
	DeleteEvent(eventID int) error
	GetEventByDate(date time.Time) (types.Event, error)
	// Add other methods as needed
}

type PreferenceRepository interface {
	GetPreferences() (types.Preferences, error)
	UpdatePreferences(prefs types.Preferences) error
	// Add other methods as needed
}
