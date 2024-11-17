// preference_repository_test.go

package repository

import (
	"testing"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPreferenceRepositorySQLite_GetPreferences(t *testing.T) {
	// Set up in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the Preferences model
	err = db.AutoMigrate(&types.Preferences{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Create an instance of PreferenceRepositorySQLite
	repo := NewPreferenceRepositorySQLite(db)

	// Insert test data
	testPrefs := types.Preferences{
		DefaultDays: "M,T,W",
		TargetDays:  "2.5",
	}
	db.Create(&testPrefs)

	// Call GetPreferences
	prefs, err := repo.GetPreferences()
	assert.NoError(t, err)
	assert.Equal(t, testPrefs.DefaultDays, prefs.DefaultDays)
	assert.Equal(t, testPrefs.TargetDays, prefs.TargetDays)
}

func TestPreferenceRepositorySQLite_GetPreferences_NoPrefs(t *testing.T) {
	// Set up in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the Preferences model
	err = db.AutoMigrate(&types.Preferences{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Create an instance of PreferenceRepositorySQLite
	repo := NewPreferenceRepositorySQLite(db)

	// Call GetPreferences without inserting any data
	prefs, err := repo.GetPreferences()
	// Expecting an error since no preferences exist
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
	assert.Empty(t, prefs)
}

func TestPreferenceRepositorySQLite_UpdatePreferences(t *testing.T) {
	// Set up in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the Preferences model
	err = db.AutoMigrate(&types.Preferences{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Create an instance of PreferenceRepositorySQLite
	repo := NewPreferenceRepositorySQLite(db)

	// Insert initial data
	testPrefs := types.Preferences{
		DefaultDays: "M,T,W",
		TargetDays:  "2.5",
	}
	db.Create(&testPrefs)

	// Update preferences
	updatedPrefs := types.Preferences{
		ID:          testPrefs.ID, // Ensure we update the same record
		DefaultDays: "M,W,F",
		TargetDays:  "3.0",
	}
	err = repo.UpdatePreferences(updatedPrefs)
	assert.NoError(t, err)

	// Retrieve preferences to check if they are updated
	prefs, err := repo.GetPreferences()
	assert.NoError(t, err)
	assert.Equal(t, updatedPrefs.DefaultDays, prefs.DefaultDays)
	assert.Equal(t, updatedPrefs.TargetDays, prefs.TargetDays)
}

func TestPreferenceRepositorySQLite_UpdatePreferences_NewRecord(t *testing.T) {
	// Set up in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the Preferences model
	err = db.AutoMigrate(&types.Preferences{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Create an instance of PreferenceRepositorySQLite
	repo := NewPreferenceRepositorySQLite(db)

	// Update preferences when no existing record
	newPrefs := types.Preferences{
		DefaultDays: "M,W,F",
		TargetDays:  "3.0",
	}
	err = repo.UpdatePreferences(newPrefs)
	assert.NoError(t, err)

	// Retrieve preferences to check if they are created
	prefs, err := repo.GetPreferences()
	assert.NoError(t, err)
	assert.Equal(t, newPrefs.DefaultDays, prefs.DefaultDays)
	assert.Equal(t, newPrefs.TargetDays, prefs.TargetDays)
}
