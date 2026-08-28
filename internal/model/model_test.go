package model

import "testing"

func TestValidateSettingsAndClone(t *testing.T) {
	settings := Settings{Color: Color{R: 1, A: 255}, Background: Color{B: 1, A: 255}, ParticleCount: 64, BeamLayers: 3, Spread: 1, Intensity: 1}
	if err := ValidateSettings(settings); err != nil {
		t.Fatal(err)
	}
	layers := []BeamLayer{{Index: 1}}
	clone := CloneLayers(layers)
	clone[0].Index = 9
	if layers[0].Index != 1 {
		t.Fatal("clone changed source")
	}
}

func TestGestureNormalization(t *testing.T) {
	if NormalizeGesture("unknown") != GestureIdle {
		t.Fatal("unknown gesture should be idle")
	}
	if !IsSpecialGesture(GestureRing) {
		t.Fatal("ring should be special")
	}
}
