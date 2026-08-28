package render

import (
	"stagebeam/internal/model"
)

type MousePoint struct {
	X float64
	Y float64
}

func GestureFromMouse(point MousePoint, width, height float64, pressed bool) (model.Gesture, float64) {
	if width <= 0 || height <= 0 {
		return model.GestureIdle, 0
	}
	normalizedX := point.X / width
	normalizedY := point.Y / height
	if normalizedX < 0 {
		normalizedX = 0
	}
	if normalizedX > 1 {
		normalizedX = 1
	}
	if normalizedY < 0 {
		normalizedY = 0
	}
	if normalizedY > 1 {
		normalizedY = 1
	}
	strength := normalizedX*0.6 + (1-normalizedY)*0.4
	if pressed {
		return model.GestureFist, strength
	}
	if normalizedY < 0.28 {
		return model.GestureRing, strength
	}
	return model.GestureOpen, strength
}

func DemoSequence() []model.GestureEvent {
	return []model.GestureEvent{
		{Gesture: model.GestureOpen, Strength: 0.35, Sequence: 1},
		{Gesture: model.GestureFist, Strength: 0.82, Sequence: 2},
		{Gesture: model.GestureRing, Strength: 0.7, Sequence: 3},
		{Gesture: model.GestureWave, Strength: 0.55, Sequence: 4},
	}
}

func PointerLabel(point MousePoint, width, height float64) string {
	gesture, strength := GestureFromMouse(point, width, height, false)
	return string(gesture) + ":" + formatStrength(strength)
}

func formatStrength(strength float64) string {
	if strength < 0.33 {
		return "soft"
	}
	if strength < 0.66 {
		return "medium"
	}
	return "strong"
}

func MouseTrail(points []MousePoint, width, height float64) []model.GestureEvent {
	trail := make([]model.GestureEvent, 0, len(points))
	for index, point := range points {
		gesture, strength := GestureFromMouse(point, width, height, index%4 == 1)
		trail = append(trail, model.GestureEvent{Gesture: gesture, Strength: strength, Sequence: int64(index + 1), Ended: index == len(points)-1})
	}
	return trail
}

func ClampPoint(point MousePoint, width, height float64) MousePoint {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if point.X < 0 {
		point.X = 0
	}
	if point.X > width {
		point.X = width
	}
	if point.Y < 0 {
		point.Y = 0
	}
	if point.Y > height {
		point.Y = height
	}
	return point
}
