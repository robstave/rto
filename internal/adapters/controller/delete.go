package controller

import (
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
)

// DeleteEvent handles deletion of a vacation event and transforms it into a remote day
func (ctlr *RTOController) DeleteEvent(c echo.Context) error {
	// Get the event ID from the URL parameter
	idParam := c.Param("id")
	eventID, err := strconv.Atoi(idParam)
	if err != nil {
		ctlr.logger.Error("Invalid event ID", "id", idParam, "error", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid event ID.",
		})
	}

	// Call the service to transform the vacation to remote
	err = ctlr.service.TransformVacationToRemote(eventID)
	if err != nil {
		ctlr.logger.Error("Error transforming vacation to remote", "eventID", eventID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Vacation day transformed into a remote day successfully.",
	})
}
