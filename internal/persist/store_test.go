package persist

import (
	"path/filepath"
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "show.db")
	settings := config.DefaultSettings()
	show := model.Show{ID: "persisted", Name: "Reopen", Settings: settings, Gesture: model.GestureIdle, Active: true, Layers: []model.BeamLayer{{Index: 0}}}
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveShow(show); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	loaded, err := reopened.LoadShow(show.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Name != show.Name || loaded.Settings.BeamLayers != show.Settings.BeamLayers {
		t.Fatalf("data did not survive reopen: %#v", loaded)
	}
}

func TestMissingShow(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "missing.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if _, err := store.LoadShow("missing"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestStoreHealthAndList(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "health.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.Health(); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveShow(model.Show{ID: "one", Name: "One", Settings: config.DefaultSettings(), Gesture: model.GestureIdle}); err != nil {
		t.Fatal(err)
	}
	shows, err := store.ListShows()
	if err != nil || len(shows) != 1 {
		t.Fatalf("unexpected show list: %#v %v", shows, err)
	}
}
