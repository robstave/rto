//go:generate mockery --name RTOBLL
package domain

import (
	"log/slog"
	"time"

	repository "github.com/robstave/rto/internal/adapters/repositories"
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
	DeleteEvent(eventID int) error
	GetEventByID(eventID int) (types.Event, error)
	TransformVacationToRemote(eventID int) error
	GetEventByDateAndType(date time.Time, eventType string) (*types.Event, error)
	GetEventsByDate(date time.Time) ([]types.Event, error)
	ClearEventsForDate(date time.Time) error

	UpdateEvent(event types.Event) error
	BulkAddEvents(events []types.Event) (*types.BulkAddResponse, error)
}

type Service struct {
	preferences    types.Preferences
	logger         *slog.Logger
	eventRepo      repository.EventRepository
	preferenceRepo repository.PreferenceRepository
}

func NewService(
	logger *slog.Logger,
	eventRepo repository.EventRepository,
	preferenceRepo repository.PreferenceRepository,
) RTOBLL {

	service := Service{
		logger:         logger,
		eventRepo:      eventRepo,
		preferenceRepo: preferenceRepo,
	}

	service.preferences = initializePreferences(&service)

	return &service
}
