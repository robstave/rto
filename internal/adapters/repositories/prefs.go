package repository

import (
	"github.com/robstave/rto/internal/domain/types"
)

func (r *PreferenceRepositorySQLite) GetPreferences() (types.Preferences, error) {
	var prefs types.Preferences
	result := r.db.First(&prefs)
	return prefs, result.Error
}

func (r *PreferenceRepositorySQLite) UpdatePreferences(prefs types.Preferences) error {
	result := r.db.Save(&prefs)
	return result.Error
}
