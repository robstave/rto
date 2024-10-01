// internal/adapters/controller/controller_test.go

package controller

import (
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/mock"
)

// MockRTOBLL is a mock implementation of the RTOBLL interface
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

func (m *MockRTOBLL) AddEvent(event types.Event) error {
	args := m.Called(event)
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
