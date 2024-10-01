package controller

import (
	"net/http"

	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/robstave/rto/internal/domain/types"

	"gorm.io/gorm"
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

func initializeDefaultPreferences(db *gorm.DB, logger *slog.Logger) error {
	var count int64
	if err := db.Model(&types.Preferences{}).Count(&count).Error; err != nil {
		logger.Error("Failed to count preferences", "error", err)
		return err
	}

	if count == 0 {
		// No preferences found; create default
		prefs := types.Preferences{
			DefaultDays: "M,T,W,Th", // Default to first 4 days in week
			TargetDays:  "2.5",
		}
		if err := db.Create(&prefs).Error; err != nil {
			logger.Error("Failed to create default preferences", "error", err)
			return err
		}
		logger.Info("Default preferences created")
	} else {
		logger.Info("Preferences already exist")
	}

	return nil
}
