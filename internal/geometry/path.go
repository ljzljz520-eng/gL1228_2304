package geometry

import (
	"math"
	"stagebeam/internal/model"
)

type Path struct {
	Points []model.Vec3
}

func NewPath(points []model.Vec3) Path {
	return Path{Points: append([]model.Vec3(nil), points...)}
}

func (path Path) At(amount float64) model.Vec3 {
	if len(path.Points) == 0 {
		return model.Vec3{}
	}
	if len(path.Points) == 1 {
		return path.Points[0]
	}
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	position := amount * float64(len(path.Points)-1)
	index := int(math.Floor(position))
	if index >= len(path.Points)-1 {
		return path.Points[len(path.Points)-1]
	}
	return Lerp(path.Points[index], path.Points[index+1], position-float64(index))
}

func (path Path) Length() float64 {
	length := 0.0
	for index := 1; index < len(path.Points); index++ {
		length += Distance(path.Points[index-1], path.Points[index])
	}
	return length
}

func BuildBeamPath(layers []model.BeamLayer, strength float64) Path {
	points := make([]model.Vec3, 0, len(layers)+1)
	points = append(points, model.Vec3{})
	for index, layer := range layers {
		point := model.Vec3{X: math.Sin(layer.Angle) * (0.3 + float64(index)*0.04), Y: math.Cos(layer.Angle) * (0.3 + float64(index)*0.04), Z: strength * (0.4 + float64(index)*0.1)}
		points = append(points, point)
	}
	return NewPath(points)
}
