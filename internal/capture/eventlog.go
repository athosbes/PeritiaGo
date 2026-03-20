package capture

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// GetSoftwareEvents parses the Windows Event Log for installation and uninstallation events in the last 180 days.
func GetSoftwareEvents() []models.TimelineEvent {
	var events []models.TimelineEvent

	// XML Query for MsiInstaller events in the last 180 days (15552000 seconds)
	// Events 1033, 11707 (Install) and 1034, 11724 (Uninstall) are common.
	query := `*[System[Provider[@Name='MsiInstaller'] and TimeCreated[timediff(@SystemTime) <= 15552000000]]]`

	out, err := exec.Command("wevtutil", "qe", "Application", "/q:"+query, "/f:text").Output()
	if err != nil {
		log.Printf("[Warning] Failed to query Event Log: %v", err)
		return events
	}

	// Parsing the text output of wevtutil
	// Each event starts with "Event[" or similar depending on the locale, but usually wevtutil /f:text follows a pattern.
	raw := string(out)
	blocks := strings.Split(raw, "\r\n\r\n")

	for _, block := range blocks {
		if !strings.Contains(block, "MsiInstaller") {
			continue
		}

		lines := strings.Split(block, "\n")
		var timestamp time.Time
		var message string
		var eventID string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Date:") {
				tsStr := strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
				// wevtutil Date format is often "2024-03-13T12:00:00.000"
				t, err := time.Parse("2006-01-02T15:04:05.000", tsStr)
				if err == nil {
					timestamp = t
				}
			} else if strings.HasPrefix(line, "Event ID:") {
				eventID = strings.TrimSpace(strings.TrimPrefix(line, "Event ID:"))
			} else if strings.HasPrefix(line, "Message:") {
				message = strings.TrimSpace(strings.TrimPrefix(line, "Message:"))
			}
		}

		if !timestamp.IsZero() && message != "" {
			eventType := "Software Event"
			if eventID == "1033" || eventID == "11707" {
				eventType = "Software Installed"
			} else if eventID == "1034" || eventID == "11724" {
				eventType = "Software Uninstalled"
			}

			events = append(events, models.TimelineEvent{
				Timestamp:   timestamp,
				Event:       eventType,
				Source:      "EventLog: Application/MsiInstaller",
				Description: fmt.Sprintf("EventID: %s | %s", eventID, message),
			})
		}
	}

	return events
}
