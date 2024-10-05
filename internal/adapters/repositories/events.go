package repository

import (
	"time"

	"github.com/robstave/rto/internal/domain/types"
)

func (r *EventRepositorySQLite) GetAllEvents() ([]types.Event, error) {
	var events []types.Event
	result := r.db.Order("date ASC").Find(&events)
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

func (r *EventRepositorySQLite) GetEventByID(eventID int) (types.Event, error) {
	var event types.Event
	result := r.db.First(&event, eventID)
	return event, result.Error
}

func (r *EventRepositorySQLite) GetEventsByType(eventType string) ([]types.Event, error) {
	var events []types.Event
	result := r.db.Where("type = ?", eventType).Order("date ASC").Find(&events)
	return events, result.Error
}

func (r *EventRepositorySQLite) GetEventsByDate(date time.Time) ([]types.Event, error) {
	var events []types.Event
	result := r.db.Where("date = ?", date).Order("date ASC").Find(&events)
	return events, result.Error
}

func (r *EventRepositorySQLite) GetEventByDateAndType(date time.Time, eventType string) (types.Event, error) {
	var event types.Event
	result := r.db.Where("date = ? AND type = ?", date, eventType).First(&event)
	return event, result.Error
}
