package stage

import (
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"testing"
)

func TestGestureKeepsBeamLayers(t *testing.T) {
	settings := config.DefaultSettings()
	settings.BeamLayers = 5
	composer := NewLayerComposer()
	layers := composer.Build(settings, model.GestureFist, 0.9)
	finished := composer.ComposeLayers(layers, model.GestureFist, true)
	if len(finished) != settings.BeamLayers {
		t.Fatalf("finished fist should retain %d layers, got %d", settings.BeamLayers, len(finished))
	}
}

func TestGestureEngineLifecycle(t *testing.T) {
	engine := NewGestureEngine()
	event := engine.Begin(model.GestureOpen, 1.2)
	if event.Strength != 1 || event.Gesture != model.GestureOpen {
		t.Fatalf("unexpected begin event: %#v", event)
	}
	ended := engine.End()
	if !ended.Ended || engineCurrent(engine) != model.GestureIdle {
		t.Fatal("engine did not end gesture")
	}
}

func TestTimelineCues(t *testing.T) {
	timeline := NewTimeline([]Cue{{Frame: 0, Gesture: model.GestureOpen, Strength: 0.4, Duration: 4}, {Frame: 5, Gesture: model.GestureRing, Strength: 0.8, Duration: 3}})
	if _, ok := timeline.At(2); !ok {
		t.Fatal("expected active cue")
	}
	if cue, ok := timeline.Next(2); !ok || cue.Gesture != model.GestureRing {
		t.Fatal("expected next ring cue")
	}
}

func engineCurrent(engine *GestureEngine) model.Gesture {
	gesture, _ := engine.Current()
	return gesture
}
