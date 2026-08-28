package stage

import (
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"testing"
)

func TestComposerSpecialGestures(t *testing.T) {
	settings := config.DefaultSettings()
	composer := NewLayerComposer()
	layers := composer.Build(settings, model.GestureRing, 0.7)
	rings := composer.ComposeLayers(layers, model.GestureRing, false)
	if len(rings) != settings.BeamLayers+3 {
		t.Fatalf("expected rings, got %d", len(rings))
	}
	fan := FanLayers(settings, 0.8)
	if len(fan) != settings.BeamLayers {
		t.Fatalf("expected fan layers, got %d", len(fan))
	}
}
