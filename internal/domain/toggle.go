package domain

import (
	"errors"
	"strconv"
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"
)

func (s *Service) ToggleAttendance(eventDate time.Time) (string, error) {
	// Retrieve all events
	events, err := s.eventRepo.GetAllEvents()
	if err != nil {
		s.logger.Error("Error retrieving events", "error", err)
		return "", err
	}

	// Find the attendance event on the given date
	found := false
	var newStatus string
	var eventToUpdate types.Event

	for _, event := range events {
		if utils.SameDay(event.Date, eventDate) && event.Type == "attendance" {
			// Toggle the IsInOffice flag
			event.IsInOffice = !event.IsInOffice
			eventToUpdate = event
			if event.IsInOffice {
				newStatus = "in"
			} else {
				newStatus = "remote"
			}
			found = true

			s.logger.Info("xmxmxmx Toggle", "Date", event.Date, "ID", event.ID)

			break
		}
	}

	if !found {
		return "", errors.New("attendance event not found on the specified date")
	}

	// Update the event in the database
	err = s.eventRepo.UpdateEvent(eventToUpdate)
	if err != nil {
		s.logger.Error("Error updating event", "error", err)
		return "", err
	}

	return newStatus, nil
}

// CalculateAttendanceStats calculates all the stats
func (s *Service) CalculateAttendanceStats() (*types.AttendanceStats, error) {
	currentYear := time.Now().Year()
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.UTC)

	// normalize ?

	allTheEvents, err := s.eventRepo.GetAllEvents()
	if err != nil {
		s.logger.Error("Error fetching  events", "error", err)
		return nil, err
	}
	inOfficeCount, totalDays := utils.CalculateInOfficeAverage(allTheEvents, startDate, endDate)

	average := 0.0
	averageDays := 0.0
	if totalDays > 0 {
		average = (float64(inOfficeCount) / float64(totalDays)) * 100
		averageDays = (float64(inOfficeCount) / float64(totalDays)) * 7 // Average days/week
	}

	// Fetch targetDays from preferences
	targetDaysStr := s.preferences.TargetDays
	targetDays, err := strconv.ParseFloat(targetDaysStr, 64)
	if err != nil {
		// Fallback to default target if parsing fails
		targetDays = 2.5
	}

	// Calculate Average Percent
	averagePercent := 0.0
	if targetDays > 0 {
		averagePercent = (averageDays / targetDays) * 100
	}

	return &types.AttendanceStats{
		InOfficeCount:  inOfficeCount,
		TotalDays:      totalDays,
		Average:        average,
		AverageDays:    averageDays,
		TargetDays:     targetDays,
		AveragePercent: averagePercent,
	}, nil
}
