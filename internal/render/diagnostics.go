package render

import (
	"fmt"
	"stagebeam/internal/model"
)

type Diagnostics struct {
	FrameID        int64
	LayerEnergy    float64
	ParticleEnergy float64
	BrightLayers   int
	DimLayers      int
	Warnings       []string
}

func Inspect(frame model.RenderFrame) Diagnostics {
	diagnostics := Diagnostics{FrameID: frame.Frame, Warnings: make([]string, 0)}
	for _, layer := range frame.Layers {
		diagnostics.LayerEnergy += layer.Energy()
		if layer.Brightness >= 0.8 {
			diagnostics.BrightLayers++
		} else {
			diagnostics.DimLayers++
		}
		if layer.Width > 0.6 {
			diagnostics.Warnings = append(diagnostics.Warnings, fmt.Sprintf("layer %d is wide", layer.Index))
		}
	}
	for _, particle := range frame.Particles {
		diagnostics.ParticleEnergy += particle.Energy
		if particle.Energy < 0.2 {
			diagnostics.Warnings = append(diagnostics.Warnings, fmt.Sprintf("particle %d needs recharge", particle.ID))
		}
	}
	if len(frame.Layers) == 0 {
		diagnostics.Warnings = append(diagnostics.Warnings, "frame has no beam layers")
	}
	if len(frame.Particles) == 0 {
		diagnostics.Warnings = append(diagnostics.Warnings, "frame has no particles")
	}
	return diagnostics
}

func (diagnostics Diagnostics) Healthy() bool {
	return len(diagnostics.Warnings) == 0 && diagnostics.LayerEnergy > 0 && diagnostics.ParticleEnergy > 0
}

func (diagnostics Diagnostics) Summary() string {
	return fmt.Sprintf("frame %d: %.2f beam energy, %.2f particle energy, %d warnings", diagnostics.FrameID, diagnostics.LayerEnergy, diagnostics.ParticleEnergy, len(diagnostics.Warnings))
}

func (diagnostics Diagnostics) WarningCount() int {
	return len(diagnostics.Warnings)
}
