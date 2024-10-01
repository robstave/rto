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
	logger *slog.Logger,

) *RTOController {

	db, err := gorm.Open(sqlite.Open("rto_attendance.db"), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		panic("Failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&types.Event{}, &types.Preferences{})

	// Initialize repositories
	eventRepo := repo.NewEventRepositorySQLite(db)
	preferenceRepo := repo.NewPreferenceRepositorySQLite(db)

	service := domain.NewService(
		logger,
		eventRepo,
		preferenceRepo,
	)

	return &RTOController{service, logger}
}
