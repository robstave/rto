package domain

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/robstave/rto/internal/domain/types"
)

func (s *Service) UpdatePreferences(defaultDays string, targetDays string) error {

	s.preferences.DefaultDays = defaultDays
	s.preferences.TargetDays = targetDays

	// Save preferences to JSON file
	preferencesFilePath := "data/preferences.json"
	if err := s.SavePreferences(preferencesFilePath); err != nil {
		s.logger.Error("Error saving preferences", "error", err)
		return err
	}
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
func initializePreferences(s *Service, filePath string) types.Preferences {
	prefs, err := loadPreferences(s, filePath)
	if err != nil {
		s.logger.Error("Error loading preferences", "error", err)

		// Set default preferences if loading fails

		return types.Preferences{
			DefaultDays: "T,W,Th,F", // Default to Tuesday, Wednesday, Thursday, Friday
			TargetDays:  "2.5",
		}

	} else {
		s.logger.Info("Preferences loaded successfully.")
		return prefs
	}
}

// LoadPreferences loads preferences from a JSON file
func loadPreferences(s *Service, filePath string) (types.Preferences, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return types.Preferences{}, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return types.Preferences{}, err
	}

	var prefs types.Preferences
	if err := json.Unmarshal(byteValue, &prefs); err != nil {
		return types.Preferences{}, err
	}

	s.preferences = prefs

	return types.Preferences{}, err
}
