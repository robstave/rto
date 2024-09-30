package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ShowPrefs renders the preferences page with current default in-office days and target
func (ctlr *RTOController) ShowPrefs(c echo.Context) error {
	//preferencesLock.RLock()
	//defer preferencesLock.RUnlock()

	data := map[string]interface{}{
		"Preferences": ctlr.service.GetPrefs(),
	}

	return c.Render(http.StatusOK, "prefs.html", data)
}

func (ctlr *RTOController) UpdatePreferences(c echo.Context) error {
	newDefaultDays := c.FormValue("defaultDays")
	newTargetDays := c.FormValue("targetDays")

	if newDefaultDays == "" || newTargetDays == "" {
		return c.String(http.StatusBadRequest, "Default Days and Target Days are required.")
	}

	// Call domain service to update preferences
	err := ctlr.service.UpdatePreferences(newDefaultDays, newTargetDays)
	if err != nil {
		ctlr.logger.Error("Error updating preferences", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to update preferences.")
	}

	return c.Redirect(http.StatusSeeOther, "/prefs")
}
