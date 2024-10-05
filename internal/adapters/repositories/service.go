//go:generate mockery --name RTOBLL

package repository

import (
	"log/slog"
)

type Service struct {
	eventRepo      EventRepository
	preferenceRepo PreferenceRepository
	logger         *slog.Logger
}
