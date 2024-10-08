// internal/adapters/controller/export.go

package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// ExportEventsMarkdown handles exporting all events as a Markdown list
func (ctlr *RTOController) ExportEventsMarkdown(c echo.Context) error {
	// Fetch all events from the service
	events := ctlr.service.GetAllEvents()

	if len(events) == 0 {
		return c.String(http.StatusOK, "No events available to export.")
	}

	// Build the Markdown content
	var sb strings.Builder
	// Add Export Date
	exportDate := time.Now().Format("January 2, 2006 at 3:04 PM")
	sb.WriteString(fmt.Sprintf("**Exported on:** %s\n\n", exportDate))

	sb.WriteString("# RTO Attendance Tracker - Events Export\n\n")
	sb.WriteString("## Events List\n\n")
	sb.WriteString("| Date | Type | Description | In Office |\n")
	sb.WriteString("| ---- | ---- | ----------- | --------- |\n")

	for _, event := range events {
		date := event.Date.Format("2006-01-02")
		eventType := strings.Title(event.Type)
		description := event.Description
		inOffice := "N/A"

		if event.Type == "attendance" {
			if event.IsInOffice {
				inOffice = "Yes"
			} else {
				inOffice = "No"
			}
		}

		// Escape pipe characters in description to prevent table formatting issues
		description = strings.ReplaceAll(description, "|", "\\|")

		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", date, eventType, description, inOffice))
	}

	markdownContent := sb.String()

	// Set the headers to prompt a file download
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=events_export.md")
	c.Response().Header().Set(echo.HeaderContentType, "text/markdown")
	return c.String(http.StatusOK, markdownContent)
}
