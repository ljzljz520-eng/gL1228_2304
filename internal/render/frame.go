package render

import (
	"fmt"
	"stagebeam/internal/model"
)

type Renderer struct {
	maxParticles int
}

func NewRenderer(maxParticles int) Renderer {
	if maxParticles < 1 {
		maxParticles = 1
	}
	return Renderer{maxParticles: maxParticles}
}

func (renderer Renderer) Frame(show model.Show) model.RenderFrame {
	particles := model.CloneParticles(show.Particles)
	if len(particles) > renderer.maxParticles {
		particles = particles[:renderer.maxParticles]
	}
	message := gestureMessage(show.Gesture, len(show.Layers))
	return model.RenderFrame{ShowID: show.ID, Frame: show.Frame, Gesture: show.Gesture, Background: show.Settings.Background, Layers: model.CloneLayers(show.Layers), Particles: particles, Message: message}
}

func gestureMessage(gesture model.Gesture, layers int) string {
	switch gesture {
	case model.GestureFist:
		return fmt.Sprintf("光柱锁定 · %d 层", layers)
	case model.GestureOpen:
		return fmt.Sprintf("扇形展开 · %d 层", layers)
	case model.GestureRing:
		return fmt.Sprintf("多层光环 · %d 层", layers)
	case model.GestureWave:
		return fmt.Sprintf("波动扫光 · %d 层", layers)
	default:
		return "等待手势"
	}
}

func BlendColor(base model.Color, overlay model.Color, amount float64) model.Color {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	blend := func(a, b uint8) uint8 { return uint8(float64(a) + (float64(b)-float64(a))*amount) }
	return model.Color{R: blend(base.R, overlay.R), G: blend(base.G, overlay.G), B: blend(base.B, overlay.B), A: blend(base.A, overlay.A)}
}

func LayerEnergy(layers []model.BeamLayer) float64 {
	energy := 0.0
	for _, layer := range layers {
		energy += layer.Brightness * layer.Width
	}
	return energy
}

func ApplyBackground(frame model.RenderFrame, background model.Color) model.RenderFrame {
	frame.Background = background
	if frame.Background.IsDark() && len(frame.Layers) > 0 {
		frame.Message += " · dark stage"
	}
	return frame
}

func HighlightLayer(frame model.RenderFrame, index int, color model.Color) model.RenderFrame {
	for position, layer := range frame.Layers {
		if layer.Index == index {
			layer.Color = color
			layer.Brightness *= 1.25
			frame.Layers[position] = layer
		}
	}
	return frame
}

func NormalizeFrame(frame model.RenderFrame) model.RenderFrame {
	if frame.Frame < 0 {
		frame.Frame = 0
	}
	frame.Layers = model.EnsureLayerIndexes(frame.Layers)
	for index := range frame.Layers {
		if frame.Layers[index].Brightness > 3 {
			frame.Layers[index].Brightness = 3
		}
	}
	return frame
}

func FramePalette(frame model.RenderFrame) []model.Color {
	palette := make([]model.Color, 0, len(frame.Layers)+1)
	palette = append(palette, frame.Background)
	for _, layer := range frame.Layers {
		duplicate := false
		for _, color := range palette {
			if color == layer.Color {
				duplicate = true
				break
			}
		}
		if !duplicate {
			palette = append(palette, layer.Color)
		}
	}
	return palette
}
