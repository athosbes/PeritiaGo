package timeline

import (
	"log"
	"sort"
	"time"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// Generate builds a unified chronological timeline out of varied forensic artifacts.
func Generate(softs []models.Software, arts []models.Artifact, evidences []models.EvidenceFile) []models.TimelineEvent {
	var timeline []models.TimelineEvent

	for _, s := range softs {
		if s.InstallDate != "" && len(s.InstallDate) == 8 {
			// Basic heuristic registry dates usually 20240401 or similar
			t, err := time.Parse("20060102", s.InstallDate)
			if err == nil {
				timeline = append(timeline, models.TimelineEvent{
					Timestamp:   t,
					Event:       "Software Installation",
					Source:      "Registry",
					Description: s.DisplayName + " " + s.DisplayVersion,
				})
			}
		}
	}

	for _, a := range arts {
		// Assuming artifacts populated timestamp as RFC3339 or general formatted string
		t, err := time.Parse(time.RFC3339, a.Timestamp)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", a.Timestamp)
		}
		if err == nil {
			timeline = append(timeline, models.TimelineEvent{
				Timestamp:   t,
				Event:       a.Type,
				Source:      a.Type,
				Description: a.Name + " - " + a.Description,
			})
		}
	}

	for _, e := range evidences {
		timeline = append(timeline, models.TimelineEvent{
			Timestamp:   e.Created,
			Event:       "File Discovered",
			Source:      "Filesystem Search",
			Description: e.Path,
		})
		
		if e.Modified.After(e.Created) {
			timeline = append(timeline, models.TimelineEvent{
				Timestamp:   e.Modified,
				Event:       "File Modified",
				Source:      "Filesystem Search",
				Description: e.Path,
			})
		}
	}

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Timestamp.Before(timeline[j].Timestamp)
	})

	log.Printf("Generated %d timeline events\n", len(timeline))
	return timeline
}
