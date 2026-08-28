package stage

import (
	"fmt"
	"stagebeam/internal/model"
)

type GestureEngine struct {
	current  model.Gesture
	sequence int64
	strength float64
}

func NewGestureEngine() *GestureEngine {
	return &GestureEngine{current: model.GestureIdle}
}

func (engine *GestureEngine) Begin(gesture model.Gesture, strength float64) model.GestureEvent {
	engine.sequence++
	engine.current = model.NormalizeGesture(gesture)
	engine.strength = clampStrength(strength)
	return model.GestureEvent{Gesture: engine.current, Strength: engine.strength, Sequence: engine.sequence}
}

func (engine *GestureEngine) End() model.GestureEvent {
	engine.sequence++
	ended := model.GestureEvent{Gesture: engine.current, Strength: engine.strength, Ended: true, Sequence: engine.sequence}
	engine.current = model.GestureIdle
	engine.strength = 0
	return ended
}

func (engine *GestureEngine) Current() (model.Gesture, float64) {
	return engine.current, engine.strength
}

func (engine *GestureEngine) Apply(event model.GestureEvent) error {
	if err := model.ValidateEvent(event); err != nil {
		return err
	}
	if event.Sequence != 0 && event.Sequence < engine.sequence {
		return fmt.Errorf("gesture sequence %d is older than %d", event.Sequence, engine.sequence)
	}
	engine.sequence = event.Sequence
	engine.current = model.NormalizeGesture(event.Gesture)
	engine.strength = clampStrength(event.Strength)
	if event.Ended {
		engine.current = model.GestureIdle
	}
	return nil
}

func (engine *GestureEngine) IsActive() bool {
	return engine.current != model.GestureIdle && engine.strength > 0
}

func (engine *GestureEngine) Strength() float64 {
	return engine.strength
}

func (engine *GestureEngine) Sequence() int64 {
	return engine.sequence
}

func (engine *GestureEngine) Reset() {
	engine.current = model.GestureIdle
	engine.strength = 0
	engine.sequence = 0
}

func clampStrength(strength float64) float64 {
	if strength < 0 {
		return 0
	}
	if strength > 1 {
		return 1
	}
	return strength
}
