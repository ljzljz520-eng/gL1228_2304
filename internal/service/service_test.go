package service

import (
	"path/filepath"
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"stagebeam/internal/persist"
	"testing"
)

func TestShowServiceWorkflow(t *testing.T) {
	store, err := persist.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	showService := NewShowService(store)
	show, err := showService.CreateShow("workflow", "Workflow", config.DefaultSettings())
	if err != nil {
		t.Fatal(err)
	}
	frame, receipt, err := showService.ApplyGesture(model.GestureEvent{ShowID: show.ID, Gesture: model.GestureOpen, Strength: 0.8, Sequence: 1})
	if err != nil || receipt.LayerCount == 0 || frame.Gesture != model.GestureOpen {
		t.Fatalf("gesture workflow failed: %#v %#v %v", frame, receipt, err)
	}
	if _, err := showService.UpdateSettings(show.ID, config.DefaultSettings()); err != nil {
		t.Fatal(err)
	}
	if len(showService.AuditEntries(show.ID)) < 3 {
		t.Fatal("expected audit entries")
	}
}

func TestDemonstrateAndReset(t *testing.T) {
	showService := NewShowService(nil)
	if _, err := showService.CreateShow("memory", "Memory", config.DefaultSettings()); err != nil {
		t.Fatal(err)
	}
	frames, err := showService.Demonstrate("memory", []model.GestureEvent{{Gesture: model.GestureRing, Strength: 0.6}})
	if err != nil || len(frames) != 1 {
		t.Fatalf("demo failed: %v", err)
	}
	if _, err := showService.Reset("memory"); err != nil {
		t.Fatal(err)
	}
}
