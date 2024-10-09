package controller

import (
	"log/slog"

	repo "github.com/robstave/rto/internal/adapters/repositories"
	"github.com/robstave/rto/internal/domain"
	"github.com/robstave/rto/internal/domain/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RTOController struct {
	service domain.RTOBLL

	logger *slog.Logger
}

func NewRTOController(
	dbPath string,
	logger *slog.Logger,

) *RTOController {

	// Read DB_PATH from environment variable, set a default if not provided

	//db, err := gorm.Open(sqlite.Open("rto_attendance.db"), &gorm.Config{})
	//dbPath = dbPath + "/rto_attendance.db"
	logger.Info("creating Database", "db", dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		panic("Failed to connect to database")
	}

	// Migrate the schema
	if err := db.AutoMigrate(&types.Event{}, &types.Preferences{}); err != nil {
		logger.Error("AutoMigrate failed", "error", err)
		panic("Failed to migrate database")
	}
	// Initialize repositories
	eventRepo := repo.NewEventRepositorySQLite(db)
	preferenceRepo := repo.NewPreferenceRepositorySQLite(db)

	// Insert default Preferences if none exist
	err = initializeDefaultPreferences(db, logger)
	if err != nil {
		logger.Error("Failed to initialize default preferences", "error", err)
		panic("Failed to initialize default preferences")
	}

	// Initialize holidays
	err = initializeHolidays(db, logger)
	if err != nil {
		logger.Error("Failed to initialize holidays", "error", err)
		panic("Failed to initialize holidays")
	}

	service := domain.NewService(
		logger,
		eventRepo,
		preferenceRepo,
	)

	return &RTOController{service, logger}
}

func NewRTOControllerWithMock(dbPath string, service domain.RTOBLL) *RTOController {
	return &RTOController{service, nil} // Pass a mock logger or nil if not used in tests
}
