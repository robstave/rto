package domain

import (
	"errors"
	"testing"

	"log/slog"

	"github.com/robstave/rto/internal/adapters/repositories/mocks"
	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePreferences_Success(t *testing.T) {
	// Initialize the mock repository
	mockRepo := new(mocks.MockPreferenceRepository)

	// Define the initial preferences
	initialPrefs := types.Preferences{
		ID:          1,
		DefaultDays: "M,T,W,Th,F",
		TargetDays:  "2.5",
	}

	// Define the updated preferences
	updatedPrefs := types.Preferences{
		ID:          1,
		DefaultDays: "T,W,Th,F",
		TargetDays:  "3.0",
	}

	// Setup expectations
	mockRepo.On("GetPreferences").Return(initialPrefs, nil)
	mockRepo.On("UpdatePreferences", updatedPrefs).Return(nil)

	// Initialize the service with the mock repository
	logger := slog.New(slog.NewTextHandler(nil, nil)) // Using a simple logger
	service := Service{
		logger:         logger,
		eventRepo:      nil, // Not needed for this test
		preferenceRepo: mockRepo,
		preferences:    initialPrefs,
	}

	// Call UpdatePreferences
	err := service.UpdatePreferences(updatedPrefs.DefaultDays, updatedPrefs.TargetDays)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, updatedPrefs, service.preferences)

	// Ensure that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUpdatePreferences_GetPreferencesError(t *testing.T) {
	// Initialize the mock repository
	mockRepo := new(mocks.MockPreferenceRepository)

	// Setup expectations
	mockRepo.On("GetPreferences").Return(types.Preferences{}, errors.New("database error"))

	// Initialize the service with the mock repository
	logger := slog.New(slog.NewTextHandler(nil, nil)) // Using a simple logger
	service := Service{
		logger:         logger,
		eventRepo:      nil, // Not needed for this test
		preferenceRepo: mockRepo,
	}

	// Call UpdatePreferences
	err := service.UpdatePreferences("M,W,F", "2.0")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	// Ensure that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUpdatePreferences_UpdatePreferencesError(t *testing.T) {
	// Initialize the mock repository
	mockRepo := new(mocks.MockPreferenceRepository)

	// Define the initial preferences
	initialPrefs := types.Preferences{
		ID:          1,
		DefaultDays: "M,T,W,Th,F",
		TargetDays:  "2.5",
	}

	// Define the updated preferences
	updatedPrefs := types.Preferences{
		ID:          1,
		DefaultDays: "T,W,Th,F",
		TargetDays:  "3.0",
	}

	// Setup expectations
	mockRepo.On("GetPreferences").Return(initialPrefs, nil)
	mockRepo.On("UpdatePreferences", updatedPrefs).Return(errors.New("update failed"))

	// Initialize the service with the mock repository
	logger := slog.New(slog.NewTextHandler(nil, nil)) // Using a simple logger
	service := Service{
		logger:         logger,
		eventRepo:      nil, // Not needed for this test
		preferenceRepo: mockRepo,
		preferences:    initialPrefs,
	}

	// Call UpdatePreferences
	err := service.UpdatePreferences(updatedPrefs.DefaultDays, updatedPrefs.TargetDays)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "update failed", err.Error())

	// Ensure that the expectations were met
	mockRepo.AssertExpectations(t)
}
