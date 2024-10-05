//go:generate mockery --name PreferenceRepository
package repository

import (
	"github.com/robstave/rto/internal/domain/types"
	"gorm.io/gorm"
)

type PreferenceRepositorySQLite struct {
	db *gorm.DB
}

func NewPreferenceRepositorySQLite(db *gorm.DB) PreferenceRepository {
	return &PreferenceRepositorySQLite{db: db}
}

type PreferenceRepository interface {
	GetPreferences() (types.Preferences, error)
	UpdatePreferences(prefs types.Preferences) error
	// Add other methods as needed
}
