package render

import (
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"testing"
)

func TestRendererAndMouseDemo(t *testing.T) {
	settings := config.DefaultSettings()
	show := model.Show{ID: "s", Gesture: model.GestureFist, Settings: settings, Layers: []model.BeamLayer{{Index: 0}}, Particles: make([]model.Particle, 4)}
	frame := NewRenderer(2).Frame(show)
	if len(frame.Particles) != 2 || frame.Message == "" {
		t.Fatalf("unexpected frame: %#v", frame)
	}
	gesture, strength := GestureFromMouse(MousePoint{X: 10, Y: 10}, 100, 100, true)
	if gesture != model.GestureFist || strength <= 0 {
		t.Fatalf("unexpected mouse gesture: %s %.2f", gesture, strength)
	}
}

func TestBlendAndEnergy(t *testing.T) {
	color := BlendColor(model.Color{R: 0, A: 255}, model.Color{R: 255, A: 255}, 0.5)
	if color.R < 120 || color.R > 135 {
		t.Fatalf("unexpected blend: %#v", color)
	}
	if LayerEnergy([]model.BeamLayer{{Width: 0.2, Brightness: 2}}) != 0.4 {
		t.Fatal("unexpected layer energy")
	}
}
