package geometry

import (
	"math"
	"stagebeam/internal/model"
	"testing"
)

func TestVectorOperations(t *testing.T) {
	vector := Normalize(model.Vec3{X: 3, Y: 4})
	if math.Abs(Length(vector)-1) > 0.0001 {
		t.Fatalf("vector was not normalized: %v", vector)
	}
	rotated := RotateZ(model.Vec3{X: 1}, math.Pi/2)
	if math.Abs(rotated.Y-1) > 0.0001 {
		t.Fatalf("unexpected rotation: %v", rotated)
	}
}

func TestFanAndGestureDirection(t *testing.T) {
	point := FanPoint(1, 3, 1, 0.8)
	if point.Y <= 0 {
		t.Fatal("fan point should face forward")
	}
	direction := GestureDirection(model.GestureFist, 0.9)
	if direction.Z < 0.9 {
		t.Fatal("fist should point upward")
	}
}
