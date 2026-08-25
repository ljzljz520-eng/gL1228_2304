package config

import (
	"fmt"
	"stagebeam/internal/model"
)

func DefaultSettings() model.Settings {
	return model.Settings{
		Color:         model.Color{R: 80, G: 190, B: 255, A: 255},
		ParticleCount: 128,
		Background:    model.Color{R: 7, G: 10, B: 24, A: 255},
		BeamLayers:    4,
		Spread:        1.25,
		Intensity:     1,
	}
}

func SettingsForPreset(name string) (model.Settings, error) {
	settings := DefaultSettings()
	switch name {
	case "midnight":
		settings.Background = model.Color{R: 3, G: 5, B: 15, A: 255}
		settings.Color = model.Color{R: 90, G: 130, B: 255, A: 255}
	case "sunrise":
		settings.Background = model.Color{R: 30, G: 9, B: 18, A: 255}
		settings.Color = model.Color{R: 255, G: 100, B: 45, A: 255}
		settings.Spread = 1.7
	case "neon":
		settings.Background = model.Color{R: 2, G: 2, B: 8, A: 255}
		settings.Color = model.Color{R: 50, G: 255, B: 180, A: 255}
		settings.BeamLayers = 6
	case "whiteout":
		settings.Background = model.Color{R: 20, G: 20, B: 24, A: 255}
		settings.Color = model.Color{R: 245, G: 245, B: 255, A: 255}
		settings.Intensity = 1.35
	default:
		return model.Settings{}, fmt.Errorf("unknown preset %q", name)
	}
	return settings, nil
}

func ClampSettings(settings model.Settings) model.Settings {
	if settings.ParticleCount < 16 {
		settings.ParticleCount = 16
	}
	if settings.ParticleCount > 4096 {
		settings.ParticleCount = 4096
	}
	if settings.BeamLayers < 1 {
		settings.BeamLayers = 1
	}
	if settings.BeamLayers > 12 {
		settings.BeamLayers = 12
	}
	if settings.Spread < 0.1 {
		settings.Spread = 0.1
	}
	if settings.Spread > 8 {
		settings.Spread = 8
	}
	if settings.Intensity < 0.1 {
		settings.Intensity = 0.1
	}
	if settings.Intensity > 2 {
		settings.Intensity = 2
	}
	return settings
}

func PresetNames() []string {
	return []string{"midnight", "sunrise", "neon", "whiteout"}
}

func ResolveSettings(settings model.Settings, preset string) (model.Settings, error) {
	if preset != "" {
		return SettingsForPreset(preset)
	}
	resolved := ClampSettings(settings)
	if err := model.ValidateSettings(resolved); err != nil {
		return model.Settings{}, err
	}
	return resolved, nil
}

func BackgroundFor(settings model.Settings, gesture model.Gesture) model.Color {
	background := settings.Background
	switch gesture {
	case model.GestureFist:
		background = model.Color{R: background.R / 2, G: background.G / 2, B: background.B / 2, A: background.A}
	case model.GestureRing:
		background = model.Color{R: minByte(int(background.R) + 8), G: minByte(int(background.G) + 5), B: minByte(int(background.B) + 15), A: background.A}
	case model.GestureWave:
		background = model.Color{R: minByte(int(background.R) + 3), G: minByte(int(background.G) + 10), B: minByte(int(background.B) + 7), A: background.A}
	}
	return background
}

func minByte(value int) uint8 {
	if value > 255 {
		return 255
	}
	if value < 0 {
		return 0
	}
	return uint8(value)
}
