package stage

import (
	"stagebeam/internal/model"
)

type Cue struct {
	Frame    int64
	Gesture  model.Gesture
	Strength float64
	Duration int64
}

type Timeline struct {
	cues []Cue
}

func NewTimeline(cues []Cue) Timeline {
	copyCues := append([]Cue(nil), cues...)
	return Timeline{cues: copyCues}
}

func (timeline Timeline) At(frame int64) (Cue, bool) {
	for _, cue := range timeline.cues {
		if frame >= cue.Frame && frame < cue.Frame+cue.Duration {
			return cue, true
		}
	}
	return Cue{}, false
}

func (timeline Timeline) Next(frame int64) (Cue, bool) {
	var next Cue
	found := false
	for _, cue := range timeline.cues {
		if cue.Frame > frame && (!found || cue.Frame < next.Frame) {
			next = cue
			found = true
		}
	}
	return next, found
}

func (timeline Timeline) Length() int64 {
	length := int64(0)
	for _, cue := range timeline.cues {
		end := cue.Frame + cue.Duration
		if end > length {
			length = end
		}
	}
	return length
}

func (timeline Timeline) Events(showID string) []model.GestureEvent {
	events := make([]model.GestureEvent, 0, len(timeline.cues))
	for index, cue := range timeline.cues {
		events = append(events, model.GestureEvent{ShowID: showID, Gesture: cue.Gesture, Strength: cue.Strength, Sequence: int64(index + 1)})
	}
	return events
}
