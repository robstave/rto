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

// TestProcessRawHolidays tests the processRawHolidays function.
func TestProcessRawHolidays(t *testing.T) {
	tests := []struct {
		name           string
		rawEvents      []RawHoliday
		expectedEvents []Event
		expectedErrors int
	}{
		{
			name: "All Valid Holidays",
			rawEvents: []RawHoliday{
				{Date: "2023-12-25", Name: "Christmas Day", Type: "holiday"},
				{Date: "2023-11-23", Name: "Thanksgiving", Type: "holiday"},
			},
			expectedEvents: []Event{
				{
					Date:        time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
					Description: "Christmas Day",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2023, time.November, 23, 0, 0, 0, 0, time.UTC),
					Description: "Thanksgiving",
					IsInOffice:  false,
					Type:        "holiday",
				},
			},
			expectedErrors: 0,
		},
		{
			name: "Some Invalid Dates",
			rawEvents: []RawHoliday{
				{Date: "2023-12-25", Name: "Christmas Day", Type: "holiday"},
				{Date: "invalid-date", Name: "Invalid Holiday", Type: "holiday"},
				{Date: "2023-11-23", Name: "Thanksgiving", Type: "holiday"},
				{Date: "2023-02-30", Name: "Impossible Date", Type: "holiday"}, // Invalid date
			},
			expectedEvents: []Event{
				{
					Date:        time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
					Description: "Christmas Day",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2023, time.November, 23, 0, 0, 0, 0, time.UTC),
					Description: "Thanksgiving",
					IsInOffice:  false,
					Type:        "holiday",
				},
			},
			expectedErrors: 2,
		},
		{
			name: "All Invalid Dates",
			rawEvents: []RawHoliday{
				{Date: "invalid-date-1", Name: "Invalid Holiday 1", Type: "holiday"},
				{Date: "invalid-date-2", Name: "Invalid Holiday 2", Type: "holiday"},
			},
			expectedEvents: []Event{},
			expectedErrors: 2,
		},
		{
			name: "Mixed Event Types with Valid and Invalid Dates",
			rawEvents: []RawHoliday{
				{Date: "2023-12-25", Name: "Christmas Day", Type: "holiday"},
				{Date: "2023-11-23", Name: "Thanksgiving", Type: "holiday"},
				{Date: "2023-01-01", Name: "New Year's Day", Type: "holiday"},
				{Date: "2023-02-30", Name: "Invalid Date", Type: "holiday"}, // Invalid
				{Date: "invalid-date", Name: "Another Invalid", Type: "holiday"},
			},
			expectedEvents: []Event{
				{
					Date:        time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
					Description: "Christmas Day",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2023, time.November, 23, 0, 0, 0, 0, time.UTC),
					Description: "Thanksgiving",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
					Description: "New Year's Day",
					IsInOffice:  false,
					Type:        "holiday",
				},
			},
			expectedErrors: 2,
		},
		{
			name:           "Empty RawEvents",
			rawEvents:      []RawHoliday{},
			expectedEvents: []Event{},
			expectedErrors: 0,
		},
		{
			name: "Duplicate Holidays on Same Date",
			rawEvents: []RawHoliday{
				{Date: "2023-12-25", Name: "Christmas Day", Type: "holiday"},
				{Date: "2023-12-25", Name: "Duplicate Christmas", Type: "holiday"},
			},
			expectedEvents: []Event{
				{
					Date:        time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
					Description: "Christmas Day",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
					Description: "Duplicate Christmas",
					IsInOffice:  false,
					Type:        "holiday",
				},
			},
			expectedErrors: 0,
		},
		{
			name: "Holidays with Different Types",
			rawEvents: []RawHoliday{
				{Date: "2023-07-04", Name: "Independence Day", Type: "holiday"},
				{Date: "2023-12-26", Name: "Boxing Day", Type: "holiday"},
				{Date: "2023-08-15", Name: "Assumption Day", Type: "holiday"},
			},
			expectedEvents: []Event{
				{
					Date:        time.Date(2023, time.July, 4, 0, 0, 0, 0, time.UTC),
					Description: "Independence Day",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2023, time.December, 26, 0, 0, 0, 0, time.UTC),
					Description: "Boxing Day",
					IsInOffice:  false,
					Type:        "holiday",
				},
				{
					Date:        time.Date(2023, time.August, 15, 0, 0, 0, 0, time.UTC),
					Description: "Assumption Day",
					IsInOffice:  false,
					Type:        "holiday",
				},
			},
			expectedErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, errorsList := processRawHolidays(tt.rawEvents)

			// Check the number of events
			if len(events) != len(tt.expectedEvents) {
				t.Errorf("Expected %d events, got %d", len(tt.expectedEvents), len(events))
			}

			// Check each event's content
			for i, expectedEvent := range tt.expectedEvents {
				if i >= len(events) {
					break
				}
				actualEvent := events[i]
				if !sameDay(actualEvent.Date, expectedEvent.Date) {
					t.Errorf("Event %d: expected date %v, got %v", i, expectedEvent.Date, actualEvent.Date)
				}
				if actualEvent.Description != expectedEvent.Description {
					t.Errorf("Event %d: expected description '%s', got '%s'", i, expectedEvent.Description, actualEvent.Description)
				}
				if actualEvent.IsInOffice != expectedEvent.IsInOffice {
					t.Errorf("Event %d: expected IsInOffice %v, got %v", i, expectedEvent.IsInOffice, actualEvent.IsInOffice)
				}
				if actualEvent.Type != expectedEvent.Type {
					t.Errorf("Event %d: expected Type '%s', got '%s'", i, expectedEvent.Type, actualEvent.Type)
				}
			}

			// Check the number of errors
			if len(errorsList) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d", tt.expectedErrors, len(errorsList))
			}
		})
	}
}

// TestLoadHolidays tests the LoadHolidays function by providing mock data.
// Note: Since LoadHolidays reads from a file, you'd typically use a temporary file for testing.
// However, since the user wants to focus on making the processing function testable, we'll skip this for now.
