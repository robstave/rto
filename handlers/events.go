package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Event represents a calendar event
type Event struct {
	Date time.Time `json:"date"`
	//Date string `json:"date"`

	Description string `json:"description"`
	IsInOffice  bool   `json:"isInOffice"`
	Type        string `json:"type"` // "attendance", "holiday", "vacation"
}

// InitializeEvents loads holidays and attendance events
func InitializeEvents() {
	holidays, err := LoadHolidays("data/holidays.json")
	if err != nil {
		log.Printf("Error loading holidays: %v", err)
		allEvents = []Event{}
		return
	}

	// Load attendance events from events.json if implementing Feature 5
	attendanceEvents, err := LoadAttendanceEvents("data/events.json")
	if err != nil {
		log.Printf("Error loading attendance events: %v", err)
		attendanceEvents = []Event{}
	}

	allEvents = append(holidays, attendanceEvents...)

	log.Printf("Events Loaded:len= %v", len(allEvents))

}

// LoadAttendanceEvents loads attendance events from a JSON file
func LoadAttendanceEvents(filePath string) ([]Event, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	//var rawEvents2 []Event
	var rawEvents2 []struct {
		Date        string `json:"date"`
		Description string `json:"description"`
		Type        string `json:"type"`
		IsInOffice  bool   `json:"isInOffice"`
	}

	log.Println("0")
	if err := json.Unmarshal(byteValue, &rawEvents2); err != nil {
		return nil, err
	}
	log.Println("a")
	// Ensure date parsing
	var events []Event
	for i, e := range rawEvents2 {
		log.Println("1")
		//parsedDate, err := time.Parse("2006-01-02", e.Date.Format("2006-01-02"))
		parsedDate, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			log.Printf("Invalid date format in events.json %v for event %v: %v", i, e.Description, err)
			continue
		}

		events = append(events, Event{
			Date:        parsedDate,
			Description: e.Description,
			IsInOffice:  e.IsInOffice, // Holidays and vacations override attendance
			Type:        e.Type,
		})

	}

	return events, nil
}

// LoadHolidays loads holidays and vacations from a JSON file
func LoadHolidays(filePath string) ([]Event, error) {

	log.Println("loading Holidays")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var rawEvents []struct {
		Date string `json:"date"`
		Name string `json:"name"`
		Type string `json:"type"`
	}

	if err := json.Unmarshal(byteValue, &rawEvents); err != nil {
		return nil, err
	}

	events := []Event{}
	for _, re := range rawEvents {
		date, err := time.Parse("2006-01-02", re.Date)
		if err != nil {
			log.Printf("Invalid date format in holidays.json: %v", err)
			continue
		}
		events = append(events, Event{
			Date:        date,
			Description: re.Name,
			IsInOffice:  false, // Holidays and vacations override attendance
			Type:        re.Type,
		})
	}

	return events, nil
}
