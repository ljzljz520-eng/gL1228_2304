package audit

import (
	"fmt"
	"stagebeam/internal/model"
)

type Report struct {
	ShowID       string
	EventCount   int
	LastAction   string
	GestureCount map[model.Gesture]int
}

func (log *Log) Report(showID string) Report {
	entries := log.ForShow(showID)
	report := Report{ShowID: showID, EventCount: len(entries), GestureCount: make(map[model.Gesture]int)}
	for _, entry := range entries {
		report.LastAction = entry.Action
		if entry.Action == "gesture.applied" {
			var gesture model.Gesture
			if len(entry.Detail) >= len("gesture=") {
				gesture = model.Gesture(entry.Detail[len("gesture="):])
			}
			report.GestureCount[gesture]++
		}
	}
	return report
}

func (report Report) Summary() string {
	return fmt.Sprintf("%s: %d events, last %s", report.ShowID, report.EventCount, report.LastAction)
}

func (report Report) HasAction(action string) bool {
	return report.LastAction == action || report.GestureCount[model.Gesture(action)] > 0
}

func (log *Log) ActionsSince(showID string, sequence int64) []string {
	entries := log.ForShow(showID)
	result := make([]string, 0)
	for _, entry := range entries {
		if entry.Sequence > sequence {
			result = append(result, entry.Action)
		}
	}
	return result
}
