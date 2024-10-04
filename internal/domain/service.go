package domain

import (
	"errors"
	"log/slog"
	"time"

	repository "github.com/robstave/rto/internal/adapters/repositories"
	"gorm.io/gorm"

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

// Existing Service struct and methods...

// TransformVacationToRemote transforms a vacation event into a remote attendance day
func (s *Service) TransformVacationToRemote(eventID int) error {
	// Retrieve the vacation event by ID
	event, err := s.GetEventByID(eventID)
	if err != nil {
		return err // Event not found or other error
	}

	if event.Type != "vacation" {
		return errors.New("only vacation events can be transformed into remote days")
	}

	// Delete the vacation event
	err = s.eventRepo.DeleteEvent(eventID)
	if err != nil {
		s.logger.Error("Failed to delete vacation event", "eventID", eventID, "error", err)
		return err
	}

	// Check if an attendance event exists on that date
	existingAttendance, err := s.eventRepo.GetEventByDateAndType(event.Date, "attendance")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No attendance event exists; create one as remote
			newAttendance := types.Event{
				Date:        event.Date,
				Description: "Remote day (transformed from vacation)",
				Type:        "attendance",
				IsInOffice:  false,
			}
			err = s.eventRepo.AddEvent(newAttendance)
			if err != nil {
				s.logger.Error("Failed to add new remote attendance event", "date", newAttendance.Date, "error", err)
				return err
			}
		} else {
			s.logger.Error("Error fetching attendance event by date", "date", event.Date, "error", err)
			return err
		}
	} else {
		// Update the existing attendance event to remote
		existingAttendance.IsInOffice = false
		existingAttendance.Description = "Remote day (transformed from vacation)"
		err = s.eventRepo.UpdateEvent(existingAttendance)
		if err != nil {
			s.logger.Error("Failed to update attendance event to remote", "eventID", existingAttendance.ID, "error", err)
			return err
		}
	}

	return nil
}
