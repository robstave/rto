package repository

import (
	"gorm.io/gorm"
)

type PreferenceRepositorySQLite struct {
	db *gorm.DB
}

func NewPreferenceRepositorySQLite(db *gorm.DB) PreferenceRepository {
	return &PreferenceRepositorySQLite{db: db}
}
