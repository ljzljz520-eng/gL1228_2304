package render

import (
	"math"
	"stagebeam/internal/model"
)

type Cell struct {
	X         int
	Y         int
	Intensity float64
	Color     model.Color
}

func Rasterize(frame model.RenderFrame, width, height int) []Cell {
	if width < 1 || height < 1 {
		return nil
	}
	cells := make([]Cell, 0, width*height/4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			intensity := cellIntensity(frame, x, y, width, height)
			if intensity > 0.02 {
				cells = append(cells, Cell{X: x, Y: y, Intensity: intensity, Color: frame.Background})
			}
		}
	}
	return cells
}

func cellIntensity(frame model.RenderFrame, x, y, width, height int) float64 {
	centerX := float64(width) / 2
	centerY := float64(height) / 2
	dx := (float64(x) - centerX) / centerX
	dy := (float64(y) - centerY) / centerY
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance > 1.2 {
		return 0
	}
	energy := 0.0
	for index, layer := range frame.Layers {
		widthFactor := 0.04 + layer.Width*0.3
		angle := math.Atan2(dy, dx) - layer.Angle
		beam := math.Exp(-(angle * angle) / (widthFactor * widthFactor))
		energy += beam * layer.Brightness / float64(index+1)
	}
	for _, particle := range frame.Particles {
		particleX := particle.Position.X
		particleY := particle.Position.Y
		distanceToParticle := math.Hypot(dx-particleX, dy-particleY)
		if distanceToParticle < particle.Size*12 {
			energy += particle.Energy * (1 - distanceToParticle/(particle.Size*12))
		}
	}
	if energy > 1 {
		return 1
	}
	return energy
}

func CellBounds(cells []Cell) (int, int, int, int) {
	if len(cells) == 0 {
		return 0, 0, 0, 0
	}
	minX, minY, maxX, maxY := cells[0].X, cells[0].Y, cells[0].X, cells[0].Y
	for _, cell := range cells[1:] {
		if cell.X < minX {
			minX = cell.X
		}
		if cell.Y < minY {
			minY = cell.Y
		}
		if cell.X > maxX {
			maxX = cell.X
		}
		if cell.Y > maxY {
			maxY = cell.Y
		}
	}
	return minX, minY, maxX, maxY
}
