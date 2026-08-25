package config

import (
	"stagebeam/internal/model"
	"testing"
)

func TestPresetsAndRequest(t *testing.T) {
	settings, err := SettingsForPreset("neon")
	if err != nil {
		t.Fatal(err)
	}
	if settings.BeamLayers != 6 {
		t.Fatalf("unexpected neon layers: %d", settings.BeamLayers)
	}
	request, err := ParseRequest([]byte(`{"name":"Demo","preset":"midnight"}`))
	if err != nil {
		t.Fatal(err)
	}
	if request.Settings.Background.R != 3 {
		t.Fatalf("unexpected background: %#v", request.Settings.Background)
	}
}

func TestClampSettings(t *testing.T) {
	settings := ClampSettings(DefaultSettings())
	if err := model.ValidateSettings(settings); err != nil {
		t.Fatal(err)
	}
}
