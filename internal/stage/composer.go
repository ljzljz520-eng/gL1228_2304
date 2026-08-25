package stage

import (
	"math"
	"stagebeam/internal/geometry"
	"stagebeam/internal/model"
)

type LayerComposer struct{}

func NewLayerComposer() LayerComposer {
	return LayerComposer{}
}

func (LayerComposer) Build(settings model.Settings, gesture model.Gesture, strength float64) []model.BeamLayer {
	count := settings.BeamLayers
	if count < 1 {
		count = 1
	}
	layers := make([]model.BeamLayer, 0, count)
	for index := 0; index < count; index++ {
		width := 0.06 + float64(index+1)*0.035
		brightness := settings.Intensity * (1 - float64(index)*0.08)
		angle := float64(index) * 0.22
		switch gesture {
		case model.GestureFist:
			width *= 0.7
			brightness *= 1.2 + strength*0.2
		case model.GestureOpen:
			width *= 1.2 + settings.Spread*0.05
			angle += settings.Spread * 0.12
		case model.GestureRing:
			width *= 1.1
			brightness *= 1.15
			angle += math.Pi / 4
		case model.GestureWave:
			angle -= 0.3
		}
		layers = append(layers, model.BeamLayer{Index: index, Width: width, Brightness: brightness, Color: settings.Color, Angle: angle})
	}
	return layers
}

func (composer LayerComposer) ComposeLayers(layers []model.BeamLayer, gesture model.Gesture, ended bool) []model.BeamLayer {
	composed := model.CloneLayers(layers)
	// Gesture state transitions (including a fist release) must retain every
	// configured beam layer, so the display never loses a halo/beam layer.
	_ = ended
	if gesture == model.GestureRing {
		composed = addRings(composed)
	}
	return composed
}

func addRings(layers []model.BeamLayer) []model.BeamLayer {
	result := make([]model.BeamLayer, 0, len(layers)+3)
	result = append(result, layers...)
	for index := 0; index < 3; index++ {
		result = append(result, model.BeamLayer{Index: len(result), Width: 0.03 + float64(index)*0.015, Brightness: 0.85 - float64(index)*0.1, Color: model.Color{R: 255, G: 200, B: 80, A: 220}, Angle: float64(index) * 2.1})
	}
	return result
}

func FanLayers(settings model.Settings, strength float64) []model.BeamLayer {
	count := settings.BeamLayers
	if count < 2 {
		count = 2
	}
	layers := make([]model.BeamLayer, count)
	for index := range layers {
		point := geometry.FanPoint(index, count, settings.Spread, strength)
		layers[index] = model.BeamLayer{Index: index, Width: 0.08 + math.Abs(point.X)*0.04, Brightness: settings.Intensity * (0.8 + point.Z*0.2), Color: settings.Color, Angle: point.X}
	}
	return layers
}

func PulseLayers(layers []model.BeamLayer, frame int64, amplitude float64) []model.BeamLayer {
	if amplitude < 0 {
		amplitude = 0
	}
	pulsed := model.CloneLayers(layers)
	phase := float64(frame%16) / 16
	for index := range pulsed {
		factor := 1 + amplitude*(0.5-phase+float64(index%3)*0.08)
		if factor < 0.2 {
			factor = 0.2
		}
		pulsed[index].Brightness *= factor
		pulsed[index].Width *= 1 + amplitude*0.1
	}
	return pulsed
}

func TintLayers(layers []model.BeamLayer, tint model.Color, amount float64) []model.BeamLayer {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	tinted := model.CloneLayers(layers)
	for index, layer := range tinted {
		layer.Color = blend(layer.Color, tint, amount)
		tinted[index] = layer
	}
	return tinted
}

func blend(a, b model.Color, amount float64) model.Color {
	interpolate := func(first, second uint8) uint8 {
		return uint8(float64(first) + (float64(second)-float64(first))*amount)
	}
	return model.Color{R: interpolate(a.R, b.R), G: interpolate(a.G, b.G), B: interpolate(a.B, b.B), A: interpolate(a.A, b.A)}
}

func LayerByIndex(layers []model.BeamLayer, index int) (model.BeamLayer, bool) {
	for _, layer := range layers {
		if layer.Index == index {
			return layer, true
		}
	}
	return model.BeamLayer{}, false
}
