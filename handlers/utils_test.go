// handlers/handlers_test.go

package handlers

import (
	"testing"
	"time"
)

// Helper function to reset global variables before each test
func resetGlobals() {
	eventsLock.Lock()
	defer eventsLock.Unlock()
	allEvents = []Event{}
}

// TestCalculateInOfficeAverage tests the calculateInOfficeAverage function with the updated signature
func TestCalculateInOfficeAverage(t *testing.T) {
	resetGlobals()

	// Define the quarter dates for testing
	currentYear := 2023
	startDate := time.Date(currentYear, time.October, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.UTC)
	totalDays := int(endDate.Sub(startDate).Hours()/24) + 1 // 92 days

	tests := []struct {
		name            string
		events          []Event
		expectedCount   int
		expectedTotal   int
		expectedAverage float64
	}{
		{
			name:            "No Events",
			events:          []Event{},
			expectedCount:   0,
			expectedTotal:   totalDays,
			expectedAverage: 0.0,
		},
		{
			name: "In-Office Events Within Quarter",
			events: []Event{
				{Date: startDate, Type: "attendance", IsInOffice: true},
				{Date: startDate.AddDate(0, 0, 1), Type: "attendance", IsInOffice: true},
				{Date: endDate, Type: "attendance", IsInOffice: true},
			},
			expectedCount:   3,
			expectedTotal:   totalDays,
			expectedAverage: (3.0 / float64(totalDays)) * 100,
		},
		{
			name: "In-Office Events Outside Quarter",
			events: []Event{
				{Date: startDate.AddDate(0, -1, 0), Type: "attendance", IsInOffice: true}, // September
				{Date: endDate.AddDate(0, 1, 0), Type: "attendance", IsInOffice: true},    // January
			},
			expectedCount:   0,
			expectedTotal:   totalDays,
			expectedAverage: 0.0,
		},
		{
			name: "Mixed Event Types",
			events: []Event{
				{Date: startDate, Type: "holiday", IsInOffice: false},
				{Date: startDate.AddDate(0, 0, 2), Type: "vacation", IsInOffice: false},
				{Date: startDate.AddDate(0, 0, 3), Type: "attendance", IsInOffice: true},
				{Date: startDate.AddDate(0, 0, 4), Type: "attendance", IsInOffice: false}, // Remote
				{Date: startDate.AddDate(0, 0, 5), Type: "attendance", IsInOffice: true},
			},
			expectedCount:   2, // Only two in-office attendance events
			expectedTotal:   totalDays,
			expectedAverage: (2.0 / float64(totalDays)) * 100,
		},
		{
			name: "Multiple In-Office Events on Same Day",
			events: []Event{
				{Date: startDate, Type: "attendance", IsInOffice: true},
				{Date: startDate, Type: "attendance", IsInOffice: true}, // Duplicate on same day
				{Date: startDate.AddDate(0, 0, 1), Type: "attendance", IsInOffice: true},
			},
			expectedCount:   3, // Assuming each event is counted separately
			expectedTotal:   totalDays,
			expectedAverage: (3.0 / float64(totalDays)) * 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the global allEvents variable

			// Call the function with parameters
			inOfficeCount, total := calculateInOfficeAverage(tt.events, startDate, endDate)

			// Verify the results
			if inOfficeCount != tt.expectedCount {
				t.Errorf("Expected inOfficeCount %d, got %d", tt.expectedCount, inOfficeCount)
			}
			if total != tt.expectedTotal {
				t.Errorf("Expected totalDays %d, got %d", tt.expectedTotal, total)
			}

			average := 0.0
			if total > 0 {
				average = (float64(inOfficeCount) / float64(total)) * 100
			}

			if average != tt.expectedAverage {
				t.Errorf("Expected average %.2f, got %.2f", tt.expectedAverage, average)
			}
		})
	}
}

// TestGetCalendarMonth tests the getCalendarMonth function
func TestGetCalendarMonth(t *testing.T) {
	tests := []struct {
		name      string
		inputDate time.Time
	}{
		{
			name:      "First Day on Sunday",
			inputDate: time.Date(2023, time.April, 2, 0, 0, 0, 0, time.UTC), // April 2, 2023 is Sunday
		},
		{
			name:      "First Day on Friday",
			inputDate: time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC), // September 1, 2023 is Friday
		},
		{
			name:      "Month with 31 Days",
			inputDate: time.Date(2023, time.March, 15, 0, 0, 0, 0, time.UTC), // March 2023 has 31 days
		},
		{
			name:      "February Non-Leap Year",
			inputDate: time.Date(2023, time.February, 10, 0, 0, 0, 0, time.UTC), // February 2023 has 28 days
		},
		{
			name:      "February Leap Year",
			inputDate: time.Date(2024, time.February, 10, 0, 0, 0, 0, time.UTC), // February 2024 has 29 days
		},
		{
			name:      "Month Starting and Ending on Same Week",
			inputDate: time.Date(2021, time.May, 15, 0, 0, 0, 0, time.UTC), // May 2021
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weeks := getCalendarMonth(tt.inputDate)

			// Basic validations
			if len(weeks) < 4 || len(weeks) > 6 {
				t.Errorf("Expected weeks between 4 and 6, got %d", len(weeks))
			}

			// Validate each week has 7 days
			for weekIdx, week := range weeks {
				if len(week) != 7 {
					t.Errorf("Week %d: expected 7 days, got %d", weekIdx+1, len(week))
				}
			}

			// Validate InMonth flags
			currentMonth := tt.inputDate.Month()
			for weekIdx, week := range weeks {
				for dayIdx, day := range week {
					expectedInMonth := day.Date.Month() == currentMonth
					if day.InMonth != expectedInMonth {
						t.Errorf("Week %d, Day %d: expected InMonth=%v, got %v",
							weekIdx+1, dayIdx+1, expectedInMonth, day.InMonth)
					}
				}
			}
		})
	}
}

// TestGetCalendarMonthDetails tests specific details of the calendar generation
func TestGetCalendarMonthDetails(t *testing.T) {
	// Example: Test March 2023
	inputDate := time.Date(2023, time.March, 15, 0, 0, 0, 0, time.UTC)
	weeks := getCalendarMonth(inputDate)

	// March 1, 2023 is a Wednesday
	// The first week should start on Sunday, February 26, 2023
	firstWeek := weeks[0]
	expectedFirstDate := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)
	if !sameDay(firstWeek[0].Date, expectedFirstDate) {
		t.Errorf("First day of first week: expected %v, got %v", expectedFirstDate, firstWeek[0].Date)
	}

	// The last week should include April 1, 2023 (Saturday)
	lastWeek := weeks[len(weeks)-1]
	expectedLastDate := time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC)
	if !sameDay(lastWeek[6].Date, expectedLastDate) {
		t.Errorf("Last day of last week: expected %v, got %v", expectedLastDate, lastWeek[6].Date)
	}

	// Validate the number of weeks for March 2023
	expectedWeeks := 5
	if len(weeks) != expectedWeeks {
		t.Errorf("Expected %d weeks for March 2023, got %d", expectedWeeks, len(weeks))
	}
}

// TestSameDay tests the sameDay function
func TestSameDay(t *testing.T) {
	tests := []struct {
		name     string
		dateA    time.Time
		dateB    time.Time
		expected bool
	}{
		{
			name:     "Identical Dates",
			dateA:    time.Date(2023, time.March, 15, 10, 30, 0, 0, time.UTC),
			dateB:    time.Date(2023, time.March, 15, 22, 45, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Different Years",
			dateA:    time.Date(2022, time.March, 15, 10, 30, 0, 0, time.UTC),
			dateB:    time.Date(2023, time.March, 15, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Different Months",
			dateA:    time.Date(2023, time.March, 15, 10, 30, 0, 0, time.UTC),
			dateB:    time.Date(2023, time.April, 15, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Different Days",
			dateA:    time.Date(2023, time.March, 14, 10, 30, 0, 0, time.UTC),
			dateB:    time.Date(2023, time.March, 15, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Same Day, Different Times",
			dateA:    time.Date(2023, time.March, 15, 0, 0, 0, 0, time.UTC),
			dateB:    time.Date(2023, time.March, 15, 23, 59, 59, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Leap Day Same",
			dateA:    time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC),
			dateB:    time.Date(2024, time.February, 29, 18, 30, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Leap Day Different",
			dateA:    time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC),
			dateB:    time.Date(2023, time.February, 28, 12, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "End of Month vs Start of Next Month",
			dateA:    time.Date(2023, time.January, 31, 23, 59, 59, 0, time.UTC),
			dateB:    time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sameDay(tt.dateA, tt.dateB)
			if result != tt.expected {
				t.Errorf("sameDay(%v, %v) = %v; want %v", tt.dateA, tt.dateB, result, tt.expected)
			}
		})
	}
}
