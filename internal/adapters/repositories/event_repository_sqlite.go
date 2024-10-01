package repository

import (
	"gorm.io/gorm"
)

type EventRepositorySQLite struct {
	db *gorm.DB
}

func NewEventRepositorySQLite(db *gorm.DB) EventRepository {
	return &EventRepositorySQLite{db: db}
}
