package controller

import (
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"path/filepath"

	"os"

	"github.com/robstave/rto/internal/domain/types"
	"github.com/robstave/rto/internal/utils"
	"gorm.io/gorm"
)

func initializeHolidays(db *gorm.DB, logger *slog.Logger) error {
	// Load holidays from JSON file
	holidaysPath := filepath.Join("static", "holidays.json")
	file, err := os.Open(holidaysPath)
	if err != nil {
		logger.Error("Failed to open holidays.json", "error", err)
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read holidays.json", "error", err)
		return err
	}

	// Define a temporary struct for unmarshaling
	type RawHoliday struct {
		Date        string `json:"date"`
		Description string `json:"description"`
		Type        string `json:"type"`       // e.g., "holiday", "vacation"
		IsInOffice  bool   `json:"isInOffice"` // Optional
	}

	var rawHolidays []RawHoliday
	if err := json.Unmarshal(byteValue, &rawHolidays); err != nil {
		logger.Error("Failed to parse holidays.json", "error", err)
		return err
	}

	// Insert holidays into the database if they don't already exist
	for _, rawHoliday := range rawHolidays {
		// Parse date
		date, err := utils.ParseDate(rawHoliday.Date)
		if err != nil {
			logger.Error("Invalid date in holidays.json", "date", rawHoliday.Date, "error", err)
			continue
		}

		// Create an Event
		holiday := types.Event{
			Date:        date,
			Description: rawHoliday.Description,
			Type:        rawHoliday.Type,
			IsInOffice:  false, // Holidays override attendance
		}

		// Check if the holiday already exists
		var count int64
		db.Model(&types.Event{}).
			Where("date = ? AND type = ?", date, holiday.Type).
			Count(&count)

		if count == 0 {
			// Insert holiday
			if err := db.Create(&holiday).Error; err != nil {
				logger.Error("Failed to insert holiday", "date", date, "error", err)
			} else {
				logger.Info("Inserted holiday", "date", date, "name", holiday.Description)
			}
		}
	}

	return nil
}
