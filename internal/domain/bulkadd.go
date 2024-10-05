package domain

import (
	"fmt"
	"strings"

	"github.com/robstave/rto/internal/domain/types"
)

// BulkAddEvents processes a list of vacation events, adding or updating them based on existing events
func (s *Service) BulkAddEvents(events []types.Event) (*types.BulkAddResponse, error) {
	// Initialize counters and result list
	var addedCount, updatedCount, skippedCount int
	var failedEvents []string
	var results []types.BulkAddResult

	s.logger.Info("---BulkAddEvents----", "events", len(events))

	for _, event := range events {
		date := event.Date
		dateStr := date.Format("2006-01-02")
		eventsOnDate, err := s.eventRepo.GetEventsByDate(date)
		if err != nil && !s.IsRecordNotFoundError(err) {
			s.logger.Error("--Error fetching events by date", "date", date, "error", err)
			failedEvents = append(failedEvents, dateStr)
			results = append(results, types.BulkAddResult{
				Date:  dateStr,
				Error: "Error fetching events for this date.",
			})
			continue
		}

		var holidayExists, vacationExists bool
		var attendanceEvent *types.Event

		for _, e := range eventsOnDate {
			switch strings.ToLower(e.Type) {
			case "holiday":
				holidayExists = true
			case "vacation":
				vacationExists = true
			case "attendance":
				attendanceEvent = &e
			}
		}

		if holidayExists {
			// Skip adding/updating if a holiday exists on the date
			skippedCount++
			results = append(results, types.BulkAddResult{
				Date:   dateStr,
				Action: "Skipped (Holiday exists)",
			})
			continue
		}

		if vacationExists {
			// Update the existing vacation event
			existingVacation, err := s.GetEventByDateAndType(date, "vacation")
			if err != nil {
				s.logger.Error("Error fetching existing vacation event", "date", date, "error", err)
				failedEvents = append(failedEvents, dateStr)
				results = append(results, types.BulkAddResult{
					Date:  dateStr,
					Error: "Error fetching existing vacation event.",
				})
				continue
			}

			existingVacation.Description = event.Description
			err = s.UpdateEvent(*existingVacation)
			if err != nil {
				s.logger.Error("Failed to update existing vacation event", "event", existingVacation, "error", err)
				failedEvents = append(failedEvents, dateStr)
				results = append(results, types.BulkAddResult{
					Date:  dateStr,
					Error: "Failed to update existing vacation event.",
				})
				continue
			}
			updatedCount++
			results = append(results, types.BulkAddResult{
				Date:        dateStr,
				Action:      "Updated existing vacation",
				Description: existingVacation.Description,
			})
			continue
		}

		if attendanceEvent != nil {
			// Update the attendance event to a vacation
			s.logger.Info("+++update", "event", event.String())

			attendanceEvent.Type = "vacation"
			attendanceEvent.Description = event.Description
			err = s.UpdateEvent(*attendanceEvent)
			if err != nil {
				s.logger.Error("Failed to update attendance event to vacation", "event", attendanceEvent, "error", err)
				failedEvents = append(failedEvents, dateStr)
				results = append(results, types.BulkAddResult{
					Date:  dateStr,
					Error: "Failed to transform attendance event to vacation.",
				})
				continue
			}
			updatedCount++
			results = append(results, types.BulkAddResult{
				Date:        dateStr,
				Action:      "Transformed attendance to vacation",
				Description: attendanceEvent.Description,
			})
			continue
		}

		// If no events exist on that date, add the vacation event
		err = s.AddEvent(event)
		if err != nil {
			s.logger.Error("Failed to add vacation event", "event", event, "error", err)
			failedEvents = append(failedEvents, dateStr)
			results = append(results, types.BulkAddResult{
				Date:  dateStr,
				Error: "Failed to add vacation event.",
			})
			continue
		}
		addedCount++
		results = append(results, types.BulkAddResult{
			Date:        dateStr,
			Action:      "Added new vacation",
			Description: event.Description,
		})
	}

	// Prepare the response message
	messageParts := []string{}
	if addedCount > 0 {
		messageParts = append(messageParts, fmt.Sprintf("Successfully added %d vacation event(s).", addedCount))
	}
	if updatedCount > 0 {
		messageParts = append(messageParts, fmt.Sprintf("Successfully updated %d event(s).", updatedCount))
	}
	if skippedCount > 0 {
		messageParts = append(messageParts, fmt.Sprintf("Skipped %d event(s) due to existing holidays.", skippedCount))
	}
	if len(failedEvents) > 0 {
		messageParts = append(messageParts, fmt.Sprintf("Failed to process events on dates: %s.", strings.Join(failedEvents, ", ")))
	}

	message := strings.Join(messageParts, " ")

	// Create the BulkAddResponse
	response := &types.BulkAddResponse{
		Success: true,
		Added:   addedCount,
		Updated: updatedCount,
		Skipped: skippedCount,
		Message: message,
		Results: results,
	}

	return response, nil
}

// IsRecordNotFoundError checks if an error is a record not found error
func (s *Service) IsRecordNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "record not found")
}
