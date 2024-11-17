// event_repository_test.go

package repository

import (
	"testing"
	"time"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Set up in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the Event model
	err = db.AutoMigrate(&types.Event{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Return a cleanup function
	return db, func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func TestEventRepositorySQLite_GetAllEvents(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	// Insert test data
	events := []types.Event{
		{
			Date:        time.Now(),
			Description: "Event 1",
			Type:        "holiday",
		},
		{
			Date:        time.Now().AddDate(0, 0, 1),
			Description: "Event 2",
			Type:        "vacation",
		},
	}
	for _, event := range events {
		db.Create(&event)
	}

	// Call GetAllEvents
	result, err := repo.GetAllEvents()
	assert.NoError(t, err)
	assert.Len(t, result, len(events))
}

func TestEventRepositorySQLite_AddEvent(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	newEvent := types.Event{
		Date:        time.Now(),
		Description: "New Event",
		Type:        "attendance",
		IsInOffice:  true,
	}

	err := repo.AddEvent(newEvent)
	assert.NoError(t, err)

	// Verify that the event was added
	var events []types.Event
	db.Find(&events)
	assert.Len(t, events, 1)
	assert.Equal(t, newEvent.Description, events[0].Description)
}

func TestEventRepositorySQLite_UpdateEvent(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	// Insert initial event
	event := types.Event{
		Date:        time.Now(),
		Description: "Original Event",
		Type:        "vacation",
	}
	db.Create(&event)

	// Update the event
	event.Description = "Updated Event"
	err := repo.UpdateEvent(event)
	assert.NoError(t, err)

	// Retrieve the event to verify the update
	var updatedEvent types.Event
	db.First(&updatedEvent, event.ID)
	assert.Equal(t, "Updated Event", updatedEvent.Description)
}

func TestEventRepositorySQLite_DeleteEvent(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	// Insert an event
	event := types.Event{
		Date:        time.Now(),
		Description: "Event to Delete",
		Type:        "attendance",
	}
	db.Create(&event)

	// Delete the event
	err := repo.DeleteEvent(int(event.ID))
	assert.NoError(t, err)

	// Verify that the event was deleted
	var count int64
	db.Model(&types.Event{}).Where("id = ?", event.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestEventRepositorySQLite_GetEventByDate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	eventDate := time.Now().Truncate(24 * time.Hour)

	// Insert an event
	event := types.Event{
		Date:        eventDate,
		Description: "Event on Specific Date",
		Type:        "holiday",
	}
	db.Create(&event)

	// Retrieve the event by date
	result, err := repo.GetEventByDate(eventDate)
	assert.NoError(t, err)
	assert.Equal(t, event.Description, result.Description)
}

func TestEventRepositorySQLite_GetEventByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	// Insert an event
	event := types.Event{
		Date:        time.Now(),
		Description: "Event by ID",
		Type:        "vacation",
	}
	db.Create(&event)

	// Retrieve the event by ID
	result, err := repo.GetEventByID(int(event.ID))
	assert.NoError(t, err)
	assert.Equal(t, event.Description, result.Description)
}

func TestEventRepositorySQLite_GetEventsByType(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	// Insert events of different types
	events := []types.Event{
		{
			Date:        time.Now(),
			Description: "Holiday Event",
			Type:        "holiday",
		},
		{
			Date:        time.Now().AddDate(0, 0, 1),
			Description: "Vacation Event",
			Type:        "vacation",
		},
		{
			Date:        time.Now().AddDate(0, 0, 2),
			Description: "Another Holiday",
			Type:        "holiday",
		},
	}
	for _, event := range events {
		db.Create(&event)
	}

	// Retrieve events by type "holiday"
	result, err := repo.GetEventsByType("holiday")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	for _, event := range result {
		assert.Equal(t, "holiday", event.Type)
	}
}

func TestEventRepositorySQLite_GetEventsByDate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	eventDate := time.Now().Truncate(24 * time.Hour)

	// Insert events on the same date
	events := []types.Event{
		{
			Date:        eventDate,
			Description: "Event 1",
			Type:        "attendance",
		},
		{
			Date:        eventDate,
			Description: "Event 2",
			Type:        "vacation",
		},
		{
			Date:        eventDate.AddDate(0, 0, 1),
			Description: "Event on Different Date",
			Type:        "holiday",
		},
	}
	for _, event := range events {
		db.Create(&event)
	}

	// Retrieve events by date
	result, err := repo.GetEventsByDate(eventDate)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	for _, event := range result {
		assert.True(t, event.Date.Equal(eventDate))
	}
}

func TestEventRepositorySQLite_GetEventByDateAndType(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	eventDate := time.Now().Truncate(24 * time.Hour)

	// Insert events
	events := []types.Event{
		{
			Date:        eventDate,
			Description: "Attendance Event",
			Type:        "attendance",
		},
		{
			Date:        eventDate,
			Description: "Vacation Event",
			Type:        "vacation",
		},
	}
	for _, event := range events {
		db.Create(&event)
	}

	// Retrieve event by date and type "vacation"
	result, err := repo.GetEventByDateAndType(eventDate, "vacation")
	assert.NoError(t, err)
	assert.Equal(t, "Vacation Event", result.Description)
}

func TestEventRepositorySQLite_GetEventByDateAndType_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewEventRepositorySQLite(db)

	eventDate := time.Now().Truncate(24 * time.Hour)

	// Insert an event of a different type
	event := types.Event{
		Date:        eventDate,
		Description: "Attendance Event",
		Type:        "attendance",
	}
	db.Create(&event)

	// Attempt to retrieve a "vacation" event on the same date
	_, err := repo.GetEventByDateAndType(eventDate, "vacation")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}
