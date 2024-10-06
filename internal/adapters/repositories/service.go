package repository

import (
	"log/slog"
)

type Service struct {
	eventRepo      EventRepository
	preferenceRepo PreferenceRepository
	logger         *slog.Logger
}
