package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// DeleteEvent handles the deletion of a vacation event
func (ctlr *RTOController) DeleteEvent(c echo.Context) error {
	// Extract the event ID from the URL parameter
	idParam := c.Param("id")
	eventID, err := strconv.Atoi(idParam)
	if err != nil {
		ctlr.logger.Error("Invalid event ID", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid event ID.",
		})
	}

	// Retrieve the event by ID
	event, err := ctlr.service.GetEventByID(eventID)
	if err != nil {
		ctlr.logger.Error("Error retrieving event", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve the event.",
		})
	}

	// Check if the event is a vacation
	if event.Type != "vacation" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"message": "Only vacation events can be deleted.",
		})
	}

	// Delete the event
	err = ctlr.service.DeleteEvent(eventID)
	if err != nil {
		ctlr.logger.Error("Error deleting event", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to delete the event.",
		})
	}

	// Successful deletion
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
