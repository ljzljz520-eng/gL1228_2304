package integration

import (
	"path/filepath"
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"stagebeam/internal/persist"
	"stagebeam/internal/service"
	"testing"
)

func TestPrimaryWorkflow(t *testing.T) {
	store, err := persist.Open(filepath.Join(t.TempDir(), "integration.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	showService := service.NewShowService(store)
	show, err := showService.CreateShow("integration", "Integration Show", config.DefaultSettings())
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := showService.ApplyGesture(model.GestureEvent{ShowID: show.ID, Gesture: model.GestureOpen, Strength: 0.5, Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	frame, err := showService.Frame(show.ID)
	if err != nil || frame.ShowID != show.ID {
		t.Fatalf("frame workflow failed: %#v %v", frame, err)
	}
}

func TestSecondaryWorkflow(t *testing.T) {
	showService := service.NewShowService(nil)
	if _, err := showService.CreateShow("secondary", "Secondary", config.DefaultSettings()); err != nil {
		t.Fatal(err)
	}
	if _, err := showService.CloseShow("secondary"); err != nil {
		t.Fatal(err)
	}
}

func TestTertiaryWorkflow(t *testing.T) {
	showService := service.NewShowService(nil)
	if _, err := showService.CreateShow("tertiary", "Tertiary", config.DefaultSettings()); err != nil {
		t.Fatal(err)
	}
	frames, err := showService.Demonstrate("tertiary", []model.GestureEvent{{Gesture: model.GestureWave, Strength: 0.4}, {Gesture: model.GestureRing, Strength: 0.9}})
	if err != nil || len(frames) != 2 {
		t.Fatalf("tertiary workflow failed: %v", err)
	}
}
