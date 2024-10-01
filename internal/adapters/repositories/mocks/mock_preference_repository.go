package mocks

import (
	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/mock"
)

// MockPreferenceRepository is a mock implementation of PreferenceRepository
type MockPreferenceRepository struct {
	mock.Mock
}

func (m *MockPreferenceRepository) GetPreferences() (types.Preferences, error) {
	args := m.Called()
	return args.Get(0).(types.Preferences), args.Error(1)
}

func (m *MockPreferenceRepository) UpdatePreferences(prefs types.Preferences) error {
	args := m.Called(prefs)
	return args.Error(0)
}
