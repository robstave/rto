//go:generate mockery --name RTOBLL
package domain

import (
	"errors"

	"github.com/robstave/rto/internal/domain/types"
	"gorm.io/gorm"
)

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
