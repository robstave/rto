package domain

import (
	"encoding/json"
	"io/ioutil"

	"github.com/robstave/rto/internal/domain/types"
)

func (s *Service) UpdatePreferences(defaultDays string, targetDays string) error {
	// Fetch current preferences from the database
	prefs, err := s.preferenceRepo.GetPreferences()
	if err != nil {
		s.logger.Error("Error fetching preferences", "error", err)
		return err
	}

	// Update preferences fields
	prefs.DefaultDays = defaultDays
	prefs.TargetDays = targetDays

	// Save updated preferences back to the database
	if err := s.preferenceRepo.UpdatePreferences(prefs); err != nil {
		s.logger.Error("Error updating preferences in repository", "error", err)
		return err
	}

	// Update the service's local copy if necessary
	s.preferences = prefs

	return nil
}

func (s *Service) SavePreferences(filePath string) error {

	data, err := json.MarshalIndent(s.preferences, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// InitializePreferences initializes preferences by loading the JSON file
func initializePreferences(s *Service) types.Preferences {
	prefs, err := s.preferenceRepo.GetPreferences()
	if err != nil {
		s.logger.Error("Error loading preferences from database", "error", err)

		// Set default preferences if loading fails or no preferences exist
		defaultPrefs := types.Preferences{
			DefaultDays: "M,T,W,Th,F", // Adjusted to include Monday by default
			TargetDays:  "2.5",
		}

		if err := s.preferenceRepo.UpdatePreferences(defaultPrefs); err != nil {
			s.logger.Error("Error setting default preferences in database", "error", err)
			// Handle as needed, possibly panic or return
		} else {
			s.logger.Info("Default preferences set in database.")
		}

		return defaultPrefs
	}

	s.logger.Info("Preferences loaded successfully from database.")
	return prefs
}
