package domain

import (
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/robstave/rto/internal/domain/types"
)

// Interface for the RTO Business logic
type RTOBLL interface {
	GetAllEvents() []types.Event
	GetPrefs() types.Preferences
	ToggleAttendance(eventDate time.Time) (string, error)
	AddEvent(event types.Event) error
	CalculateAttendanceStats() (*types.AttendanceStats, error)
	UpdatePreferences(defaultDays string, targetDays string) error
	AddDefaultDays() error
}

// Global variable to store all events and manage thread safety.
// Ideally this is in the service.  but lets park it here for now
var (
	allEvents  []types.Event
	eventsLock sync.RWMutex
)

type Service struct {
	preferences types.Preferences
	logger      *slog.Logger
}

func NewService(
	logger *slog.Logger,

) RTOBLL {

	service := Service{
		logger: logger,
	}

	preferencesPath := filepath.Join("data", "preferences.json")
	service.preferences = initializePreferences(&service, preferencesPath)
	holidaysPath := filepath.Join("data", "holidays.json")
	eventsPath := filepath.Join("data", "events.json")
	initializeEvents(&service, holidaysPath, eventsPath)

	return &service
}
