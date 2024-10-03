package mocks

import (
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/mock"
)

type MockRTOBLL struct {
	mock.Mock
}

func (m *MockRTOBLL) GetAllEvents() []types.Event {
	args := m.Called()
	return args.Get(0).([]types.Event)
}

func (m *MockRTOBLL) GetPrefs() types.Preferences {
	args := m.Called()
	return args.Get(0).(types.Preferences)
}

func (m *MockRTOBLL) ToggleAttendance(eventDate time.Time) (string, error) {
	args := m.Called(eventDate)
	return args.String(0), args.Error(1)
}

func (m *MockRTOBLL) GetEventByID(eventID int) (types.Event, error) {
	args := m.Called(eventID)
	return args.Get(0).(types.Event), args.Error(1)
}
func (m *MockRTOBLL) AddEvent(event types.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
func (m *MockRTOBLL) DeleteEvent(eventID int) error {
	args := m.Called(eventID)
	return args.Error(0)
}

func (m *MockRTOBLL) CalculateAttendanceStats() (*types.AttendanceStats, error) {
	args := m.Called()
	return args.Get(0).(*types.AttendanceStats), args.Error(1)
}

func (m *MockRTOBLL) UpdatePreferences(defaultDays string, targetDays string) error {
	args := m.Called(defaultDays, targetDays)
	return args.Error(0)
}

func (m *MockRTOBLL) AddDefaultDays() error {
	args := m.Called()
	return args.Error(0)
}
