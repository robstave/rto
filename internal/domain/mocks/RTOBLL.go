// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"

	types "github.com/robstave/rto/internal/domain/types"
)

// RTOBLL is an autogenerated mock type for the RTOBLL type
type RTOBLL struct {
	mock.Mock
}

// AddDefaultDays provides a mock function with given fields:
func (_m *RTOBLL) AddDefaultDays() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddEvent provides a mock function with given fields: event
func (_m *RTOBLL) AddEvent(event types.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BulkAddEvents provides a mock function with given fields: events
func (_m *RTOBLL) BulkAddEvents(events []types.Event) (*types.BulkAddResponse, error) {
	ret := _m.Called(events)

	var r0 *types.BulkAddResponse
	if rf, ok := ret.Get(0).(func([]types.Event) *types.BulkAddResponse); ok {
		r0 = rf(events)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.BulkAddResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]types.Event) error); ok {
		r1 = rf(events)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CalculateAttendanceStats provides a mock function with given fields:
func (_m *RTOBLL) CalculateAttendanceStats() (*types.AttendanceStats, error) {
	ret := _m.Called()

	var r0 *types.AttendanceStats
	if rf, ok := ret.Get(0).(func() *types.AttendanceStats); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.AttendanceStats)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClearEventsForDate provides a mock function with given fields: date
func (_m *RTOBLL) ClearEventsForDate(date time.Time) error {
	ret := _m.Called(date)

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Time) error); ok {
		r0 = rf(date)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteEvent provides a mock function with given fields: eventID
func (_m *RTOBLL) DeleteEvent(eventID int) error {
	ret := _m.Called(eventID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllEvents provides a mock function with given fields:
func (_m *RTOBLL) GetAllEvents() []types.Event {
	ret := _m.Called()

	var r0 []types.Event
	if rf, ok := ret.Get(0).(func() []types.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Event)
		}
	}

	return r0
}

// GetEventByDateAndType provides a mock function with given fields: date, eventType
func (_m *RTOBLL) GetEventByDateAndType(date time.Time, eventType string) (*types.Event, error) {
	ret := _m.Called(date, eventType)

	var r0 *types.Event
	if rf, ok := ret.Get(0).(func(time.Time, string) *types.Event); ok {
		r0 = rf(date, eventType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Event)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(time.Time, string) error); ok {
		r1 = rf(date, eventType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventByID provides a mock function with given fields: eventID
func (_m *RTOBLL) GetEventByID(eventID int) (types.Event, error) {
	ret := _m.Called(eventID)

	var r0 types.Event
	if rf, ok := ret.Get(0).(func(int) types.Event); ok {
		r0 = rf(eventID)
	} else {
		r0 = ret.Get(0).(types.Event)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsByDate provides a mock function with given fields: date
func (_m *RTOBLL) GetEventsByDate(date time.Time) ([]types.Event, error) {
	ret := _m.Called(date)

	var r0 []types.Event
	if rf, ok := ret.Get(0).(func(time.Time) []types.Event); ok {
		r0 = rf(date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Event)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(time.Time) error); ok {
		r1 = rf(date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPrefs provides a mock function with given fields:
func (_m *RTOBLL) GetPrefs() types.Preferences {
	ret := _m.Called()

	var r0 types.Preferences
	if rf, ok := ret.Get(0).(func() types.Preferences); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.Preferences)
	}

	return r0
}

// ToggleAttendance provides a mock function with given fields: eventDate
func (_m *RTOBLL) ToggleAttendance(eventDate time.Time) (string, error) {
	ret := _m.Called(eventDate)

	var r0 string
	if rf, ok := ret.Get(0).(func(time.Time) string); ok {
		r0 = rf(eventDate)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(time.Time) error); ok {
		r1 = rf(eventDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransformVacationToRemote provides a mock function with given fields: eventID
func (_m *RTOBLL) TransformVacationToRemote(eventID int) error {
	ret := _m.Called(eventID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateEvent provides a mock function with given fields: event
func (_m *RTOBLL) UpdateEvent(event types.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePreferences provides a mock function with given fields: defaultDays, targetDays
func (_m *RTOBLL) UpdatePreferences(defaultDays string, targetDays string) error {
	ret := _m.Called(defaultDays, targetDays)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(defaultDays, targetDays)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRTOBLL interface {
	mock.TestingT
	Cleanup(func())
}

// NewRTOBLL creates a new instance of RTOBLL. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRTOBLL(t mockConstructorTestingTNewRTOBLL) *RTOBLL {
	mock := &RTOBLL{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
