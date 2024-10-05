//go:generate mockery --name EventRepository
package repository

import (
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"gorm.io/gorm"
)

type EventRepositorySQLite struct {
	db *gorm.DB
}

type EventRepository interface {
	GetAllEvents() ([]types.Event, error)
	AddEvent(event types.Event) error
	UpdateEvent(event types.Event) error
	DeleteEvent(eventID int) error
	GetEventByDate(date time.Time) (types.Event, error)
	GetEventByID(eventID int) (types.Event, error)
	GetEventsByType(eventType string) ([]types.Event, error)

	// New method to fetch event by date and type
	GetEventByDateAndType(date time.Time, eventType string) (types.Event, error)
	GetEventsByDate(date time.Time) ([]types.Event, error)
}

func NewEventRepositorySQLite(db *gorm.DB) EventRepository {
	return &EventRepositorySQLite{db: db}
}
