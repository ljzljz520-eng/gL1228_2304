package model

import (
	"fmt"
	"strings"
)

func ValidateSettings(settings Settings) error {
	if err := ValidateColor(settings.Color); err != nil {
		return fmt.Errorf("beam color: %w", err)
	}
	if err := ValidateColor(settings.Background); err != nil {
		return fmt.Errorf("background color: %w", err)
	}
	if settings.ParticleCount < 16 || settings.ParticleCount > 4096 {
		return fmt.Errorf("particle count must be between 16 and 4096")
	}
	if settings.BeamLayers < 1 || settings.BeamLayers > 12 {
		return fmt.Errorf("beam layers must be between 1 and 12")
	}
	if settings.Spread <= 0 || settings.Spread > 8 {
		return fmt.Errorf("spread must be greater than zero and no more than eight")
	}
	if settings.Intensity < 0.1 || settings.Intensity > 2 {
		return fmt.Errorf("intensity must be between 0.1 and 2")
	}
	return nil
}

func ValidateColor(color Color) error {
	if color.A == 0 {
		return fmt.Errorf("color alpha cannot be zero")
	}
	if color.R == 0 && color.G == 0 && color.B == 0 {
		return fmt.Errorf("color cannot be black")
	}
	return nil
}

func ValidateShowName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("show name is required")
	}
	if len(trimmed) > 80 {
		return fmt.Errorf("show name exceeds 80 characters")
	}
	return nil
}

func NormalizeGesture(gesture Gesture) Gesture {
	switch gesture {
	case GestureOpen, GestureFist, GestureRing, GestureWave:
		return gesture
	default:
		return GestureIdle
	}
}

func IsSpecialGesture(gesture Gesture) bool {
	return gesture == GestureRing || gesture == GestureWave
}

func CloneLayers(layers []BeamLayer) []BeamLayer {
	if len(layers) == 0 {
		return nil
	}
	cloned := make([]BeamLayer, len(layers))
	copy(cloned, layers)
	return cloned
}

func CloneParticles(particles []Particle) []Particle {
	if len(particles) == 0 {
		return nil
	}
	cloned := make([]Particle, len(particles))
	copy(cloned, particles)
	return cloned
}

func ValidateEvent(event GestureEvent) error {
	if strings.TrimSpace(event.ShowID) == "" {
		return fmt.Errorf("show id is required")
	}
	if event.Strength < 0 || event.Strength > 1 {
		return fmt.Errorf("gesture strength must be between zero and one")
	}
	if event.Sequence < 0 {
		return fmt.Errorf("gesture sequence cannot be negative")
	}
	return nil
}

func ValidateShow(show Show) error {
	if err := ValidateShowName(show.Name); err != nil {
		return err
	}
	if err := ValidateSettings(show.Settings); err != nil {
		return err
	}
	if show.Frame < 0 {
		return fmt.Errorf("frame cannot be negative")
	}
	if show.Gesture == "" {
		return fmt.Errorf("gesture is required")
	}
	return nil
}

func EnsureLayerIndexes(layers []BeamLayer) []BeamLayer {
	result := CloneLayers(layers)
	for index := range result {
		result[index].Index = index
		if result[index].Width <= 0 {
			result[index].Width = 0.05
		}
		if result[index].Brightness < 0 {
			result[index].Brightness = 0
		}
	}
	return result
}
