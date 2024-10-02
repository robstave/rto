package repository

import (
	"time"

	"github.com/robstave/rto/internal/domain/types"
)

func (r *EventRepositorySQLite) GetAllEvents() ([]types.Event, error) {
	var events []types.Event
	result := r.db.Find(&events)
	return events, result.Error
}

func (r *EventRepositorySQLite) AddEvent(event types.Event) error {
	result := r.db.Create(&event)
	return result.Error
}

func (r *EventRepositorySQLite) UpdateEvent(event types.Event) error {
	result := r.db.Save(&event)
	return result.Error
}

func (r *EventRepositorySQLite) DeleteEvent(eventID int) error {
	result := r.db.Delete(&types.Event{}, eventID)
	return result.Error
}

func (r *EventRepositorySQLite) GetEventByDate(date time.Time) (types.Event, error) {
	var event types.Event
	result := r.db.Where("date = ?", date).First(&event)
	return event, result.Error
}