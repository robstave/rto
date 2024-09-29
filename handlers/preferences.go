package handlers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// Preferences represents user preferences
type Preferences struct {
	DefaultDays string `json:"defaultDays"` // e.g., "M,T,W,Th,F"
	TargetDays  string `json:"targetDays"`  // e.g., "2.5"

}

// Global variables to store preferences and manage thread safety
var (
	preferences     Preferences
	preferencesLock sync.RWMutex
)

// InitializePreferences initializes preferences by loading the JSON file
func InitializePreferences(filePath string) {
	err := LoadPreferences(filePath)
	if err != nil {
		logger.Error("Error loading preferences", "error", err)

		// Set default preferences if loading fails
		preferencesLock.Lock()
		preferences = Preferences{
			DefaultDays: "T,W,Th,F", // Default to Tuesday, Wednesday, Thursday, Friday
		}
		preferencesLock.Unlock()
	} else {
		logger.Info("Preferences loaded successfully.")

	}
}

// LoadPreferences loads preferences from a JSON file
func LoadPreferences(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var prefs Preferences
	if err := json.Unmarshal(byteValue, &prefs); err != nil {
		return err
	}

	// Store preferences globally with thread safety
	preferencesLock.Lock()
	preferences = prefs
	preferencesLock.Unlock()

	return nil
}
