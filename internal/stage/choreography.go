package stage

import (
	"stagebeam/internal/model"
)

type Choreography struct {
	Name     string
	Timeline Timeline
}

func DefaultChoreography() Choreography {
	return Choreography{Name: "four beat beam", Timeline: NewTimeline([]Cue{{Frame: 0, Gesture: model.GestureOpen, Strength: 0.4, Duration: 24}, {Frame: 24, Gesture: model.GestureFist, Strength: 0.9, Duration: 18}, {Frame: 42, Gesture: model.GestureRing, Strength: 0.75, Duration: 18}, {Frame: 60, Gesture: model.GestureWave, Strength: 0.55, Duration: 24}})}
}

func (choreography Choreography) EventAt(showID string, frame int64) (model.GestureEvent, bool) {
	cue, ok := choreography.Timeline.At(frame)
	if !ok {
		return model.GestureEvent{}, false
	}
	return model.GestureEvent{ShowID: showID, Gesture: cue.Gesture, Strength: cue.Strength, Sequence: frame + 1, Ended: frame+1 == cue.Frame+cue.Duration}, true
}

func (choreography Choreography) Events(showID string) []model.GestureEvent {
	return choreography.Timeline.Events(showID)
}

func (choreography Choreography) Duration() int64 {
	return choreography.Timeline.Length()
}

func (choreography Choreography) CueCount() int {
	return len(choreography.Timeline.cues)
}
