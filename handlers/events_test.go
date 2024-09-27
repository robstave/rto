package handlers

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

// TestLoadAttendanceEvents is a unit test for LoadAttendanceEvents function.
func TestLoadAttendanceEvents(t *testing.T) {
	// Prepare a sample JSON file for testing
	jsonData := `[
		{
			"date": "2024-09-25",
			"description": "Meeting",
			"type": "attendance",
			"isInOffice": true
		},
		{
			"date": "2024-09-26",
			"description": "Holiday",
			"type": "holiday",
			"isInOffice": false
		}
	]`
	tempFile, err := ioutil.TempFile("", "events.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(jsonData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the temp file to flush changes
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Load the events from the file
	events, err := LoadAttendanceEvents(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load events: %v", err)
	}

	// Validate the loaded events
	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	// Check the first event
	event1 := events[0]
	expectedDate1, _ := time.Parse("2006-01-02", "2024-09-25")
	if !event1.Date.Equal(expectedDate1) {
		t.Errorf("Expected first event date to be %v, got %v", expectedDate1, event1.Date)
	}
	if event1.Description != "Meeting" {
		t.Errorf("Expected first event description to be 'Meeting', got %s", event1.Description)
	}
	if event1.Type != "attendance" {
		t.Errorf("Expected first event type to be 'attendance', got %s", event1.Type)
	}
	if event1.IsInOffice != true {
		t.Errorf("Expected first event IsInOffice to be true, got %v", event1.IsInOffice)
	}

	// Check the second event
	event2 := events[1]
	expectedDate2, _ := time.Parse("2006-01-02", "2024-09-26")
	if !event2.Date.Equal(expectedDate2) {
		t.Errorf("Expected second event date to be %v, got %v", expectedDate2, event2.Date)
	}
	if event2.Description != "Holiday" {
		t.Errorf("Expected second event description to be 'Holiday', got %s", event2.Description)
	}
	if event2.Type != "holiday" {
		t.Errorf("Expected second event type to be 'holiday', got %s", event2.Type)
	}
	if event2.IsInOffice != false {
		t.Errorf("Expected second event IsInOffice to be false, got %v", event2.IsInOffice)
	}
}

// TestLoadHolidays tests the LoadHolidays function with various scenarios
func TestLoadHolidays(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		jsonContent    string
		expectedEvents []Event
		expectError    bool
	}{
		{
			name: "Valid Holidays and Vacations",
			jsonContent: `[
                {
                    "date": "2024-10-31",
                    "name": "Halloween",
                    "type": "holiday"
                },
                {
                    "date": "2024-11-25",
                    "name": "Thanksgiving",
                    "type": "holiday"
                },
                {
                    "date": "2024-12-24",
                    "name": "Christmas Eve",
                    "type": "vacation"
                }
            ]`,
			expectedEvents: []Event{
				{
					Date:        time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC),
					Description: "Halloween",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2024, 11, 25, 0, 0, 0, 0, time.UTC),
					Description: "Thanksgiving",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC),
					Description: "Christmas Eve",
					IsInOffice:  false,
					Type:        "vacation",
				},
			},
			expectError: false,
		},
		{
			name:           "Empty JSON Array",
			jsonContent:    `[]`,
			expectedEvents: []Event{},
			expectError:    false,
		},
		{
			name:           "Invalid JSON Format",
			jsonContent:    `[{ "date": "2024-10-31", "name": "Halloween", "type": "holiday"`, // Missing closing bracket
			expectedEvents: nil,
			expectError:    true,
		},
		{
			name: "Missing Fields",
			jsonContent: `[
                {
                    "date": "2024-10-31",
                    "type": "holiday"
                },
                {
                    "name": "Thanksgiving",
                    "type": "holiday"
                }
            ]`,
			expectedEvents: []Event{
				{
					Date:        time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC),
					Description: "", // Name is missing
					IsInOffice:  false,
					Type:        "holiday",
				},
			},
			expectError: false, // The function skips invalid entries but does not return an error
		},
		{
			name: "Invalid Date Format",
			jsonContent: `[
                {
                    "date": "31-10-2024",
                    "name": "Halloween",
                    "type": "holiday"
                },
                {
                    "date": "2024-11-25",
                    "name": "Thanksgiving",
                    "type": "holiday"
                }
            ]`,
			expectedEvents: []Event{
				{
					Date:        time.Date(2024, 11, 25, 0, 0, 0, 0, time.UTC),
					Description: "Thanksgiving",
					IsInOffice:  false,
					Type:        "holiday",
				},
			},
			expectError: false, // The function skips entries with invalid dates but does not return an error
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file
			tmpFile, err := ioutil.TempFile("", "holidays_test_*.json")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name()) // Clean up

			// Write the JSON content to the temp file
			if _, err := tmpFile.Write([]byte(tc.jsonContent)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Call LoadHolidays with the temp file path
			events, err := LoadHolidays(tmpFile.Name())

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				// If an error is expected, no need to check further
				return
			} else {
				if err != nil {
					t.Errorf("Did not expect error but got: %v", err)
					return
				}
			}

			// Compare the expected and actual events
			if !reflect.DeepEqual(events, tc.expectedEvents) {
				t.Errorf("Expected events:\n%v\nGot events:\n%v", tc.expectedEvents, events)
			}
		})
	}
}
