package domain

import (
	"errors"
	"strconv"
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"
)

func (s *Service) ToggleAttendance(eventDate time.Time) (string, error) {
	eventsLock.Lock()
	defer eventsLock.Unlock()

	// Find the attendance event on the given date
	found := false
	var newStatus string
	for i, event := range allEvents {
		if utils.SameDay(event.Date, eventDate) && event.Type == "attendance" {
			// Toggle the IsInOffice flag
			allEvents[i].IsInOffice = !event.IsInOffice
			if allEvents[i].IsInOffice {
				newStatus = "in"
			} else {
				newStatus = "remote"
			}
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("attendance event not found on the specified date")
	}

	// Save to events.json
	eventsFilePath := "data/events.json"
	if err := SaveEvents(eventsFilePath); err != nil {
		s.logger.Error("Error saving events", "error", err)
		return "", err
	}

	return newStatus, nil
}

func (s *Service) CalculateAttendanceStats() (*types.AttendanceStats, error) {
	currentYear := time.Now().Year()
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.Local)

	inOfficeCount, totalDays := utils.CalculateInOfficeAverage(allEvents, startDate, endDate)

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

	return &types.AttendanceStats{
		InOfficeCount: inOfficeCount,
		TotalDays:     totalDays,
		Average:       average,
		AverageDays:   averageDays,
		TargetDays:    targetDays,
	}, nil
}
