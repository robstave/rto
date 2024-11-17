// domain/toggle_test.go

package domain

import (
	"testing"
	"time"

	"log/slog"
	"os"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/assert"
)

// MockEventRepository is a mock implementation of EventRepository
type MockEventRepository struct {
	GetAllEventsFunc func() ([]types.Event, error)
	UpdateEventFunc  func(event types.Event) error
}

func (m *MockEventRepository) GetAllEvents() ([]types.Event, error) {
	return m.GetAllEventsFunc()
}

func (m *MockEventRepository) AddEvent(event types.Event) error {
	return nil
}

func (m *MockEventRepository) UpdateEvent(event types.Event) error {
	return m.UpdateEventFunc(event)
}

func (m *MockEventRepository) DeleteEvent(eventID int) error {
	return nil
}

func (m *MockEventRepository) GetEventByDate(date time.Time) (types.Event, error) {
	return types.Event{}, nil
}

func (m *MockEventRepository) GetEventByID(eventID int) (types.Event, error) {
	return types.Event{}, nil
}

func (m *MockEventRepository) GetEventsByType(eventType string) ([]types.Event, error) {
	return nil, nil
}

func (m *MockEventRepository) GetEventByDateAndType(date time.Time, eventType string) (types.Event, error) {
	return types.Event{}, nil
}

func (m *MockEventRepository) GetEventsByDate(date time.Time) ([]types.Event, error) {
	return nil, nil
}

// TestToggleAttendance_Success tests the successful toggling of attendance
func TestToggleAttendance_Success(t *testing.T) {
	// Set up mock data
	eventDate := time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)
	existingEvent := types.Event{
		ID:         1,
		Date:       eventDate,
		Type:       "attendance",
		IsInOffice: false,
	}

	// Create mock EventRepository
	mockEventRepo := &MockEventRepository{
		GetAllEventsFunc: func() ([]types.Event, error) {
			return []types.Event{existingEvent}, nil
		},
		UpdateEventFunc: func(event types.Event) error {
			// Verify that the event's IsInOffice was toggled
			assert.Equal(t, existingEvent.ID, event.ID)
			assert.Equal(t, !existingEvent.IsInOffice, event.IsInOffice)
			return nil
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)) // Using a simple logger
	// Create service with mock repository
	service := &Service{
		eventRepo: mockEventRepo,
		logger:    logger, // Replace with a mock logger if needed
	}

	newStatus, err := service.ToggleAttendance(eventDate)
	assert.NoError(t, err)
	assert.Equal(t, "in", newStatus)
}

// TestToggleAttendance_EventNotFound tests the case where the attendance event is not found
func TestToggleAttendance_EventNotFound(t *testing.T) {
	// Set up mock data
	eventDate := time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)

	// Create mock EventRepository that returns no events
	mockEventRepo := &MockEventRepository{
		GetAllEventsFunc: func() ([]types.Event, error) {
			return []types.Event{}, nil
		},
		UpdateEventFunc: func(event types.Event) error {
			return nil
		},
	}

	// Create service with mock repository
	service := &Service{
		eventRepo: mockEventRepo,
		logger:    nil,
	}

	newStatus, err := service.ToggleAttendance(eventDate)
	assert.Error(t, err)
	assert.Equal(t, "", newStatus)
	assert.Equal(t, "attendance event not found on the specified date", err.Error())
}

// TestCalculateAttendanceStats tests the calculation of attendance statistics
func TestCalculateAttendanceStats(t *testing.T) {
	// Set up mock data
	currentYear := 2024
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.UTC)

	events := []types.Event{
		{
			ID:         1,
			Date:       time.Date(currentYear, time.October, 2, 0, 0, 0, 0, time.UTC),
			Type:       "attendance",
			IsInOffice: true,
		},
		{
			ID:         2,
			Date:       time.Date(currentYear, time.October, 3, 0, 0, 0, 0, time.UTC),
			Type:       "attendance",
			IsInOffice: false,
		},
		{
			ID:         3,
			Date:       time.Date(currentYear, time.October, 4, 0, 0, 0, 0, time.UTC),
			Type:       "attendance",
			IsInOffice: true,
		},
		{
			ID:         4,
			Date:       time.Date(currentYear, time.November, 5, 0, 0, 0, 0, time.UTC),
			Type:       "attendance",
			IsInOffice: true,
		},
		// Event outside the date range
		{
			ID:         5,
			Date:       time.Date(currentYear, time.September, 30, 0, 0, 0, 0, time.UTC),
			Type:       "attendance",
			IsInOffice: true,
		},
	}

	// Create mock EventRepository
	mockEventRepo := &MockEventRepository{
		GetAllEventsFunc: func() ([]types.Event, error) {
			return events, nil
		},
	}

	// Create service with mock repository and preferences
	service := &Service{
		eventRepo: mockEventRepo,
		logger:    nil,
		preferences: types.Preferences{
			TargetDays: "2.5",
		},
	}

	stats, err := service.CalculateAttendanceStats()
	assert.NoError(t, err)

	// Expected in-office count is 3 (events with IsInOffice true within date range)
	expectedInOfficeCount := 3

	// Total days in the quarter is from Oct 1 to Dec 31
	expectedTotalDays := int(endDate.Sub(startDate).Hours()/24) + 1

	// Calculate expected average and averageDays
	expectedAverage := (float64(expectedInOfficeCount) / float64(expectedTotalDays)) * 100
	expectedAverageDays := (float64(expectedInOfficeCount) / float64(expectedTotalDays)) * 7

	expectedTargetDays := 2.5
	expectedAveragePercent := 0.0
	if expectedTargetDays > 0 {
		expectedAveragePercent = (expectedAverageDays / expectedTargetDays) * 100
	}

	assert.Equal(t, expectedInOfficeCount, stats.InOfficeCount)
	assert.Equal(t, expectedTotalDays, stats.TotalDays)
	assert.InEpsilon(t, expectedAverage, stats.Average, 0.0001)
	assert.InEpsilon(t, expectedAverageDays, stats.AverageDays, 0.0001)
	assert.Equal(t, expectedTargetDays, stats.TargetDays)
	assert.InEpsilon(t, expectedAveragePercent, stats.AveragePercent, 0.0001)
}
