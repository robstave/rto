package controller

import (
	"log/slog"

	"github.com/robstave/rto/internal/domain"
)

type RTOController struct {
	service domain.RTOBLL

	logger *slog.Logger
}

func NewRTOController(
	logger *slog.Logger,

) *RTOController {

	service := domain.NewService(
		logger,
	)

	return &RTOController{service, logger}
}
