package geometry

import (
	"math"
	"stagebeam/internal/model"
)

func Add(a, b model.Vec3) model.Vec3 {
	return model.Vec3{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

func Scale(vector model.Vec3, factor float64) model.Vec3 {
	return model.Vec3{X: vector.X * factor, Y: vector.Y * factor, Z: vector.Z * factor}
}

func Length(vector model.Vec3) float64 {
	return math.Sqrt(vector.X*vector.X + vector.Y*vector.Y + vector.Z*vector.Z)
}

func Normalize(vector model.Vec3) model.Vec3 {
	length := Length(vector)
	if length == 0 {
		return model.Vec3{}
	}
	return Scale(vector, 1/length)
}

func RotateZ(vector model.Vec3, radians float64) model.Vec3 {
	cosine := math.Cos(radians)
	sine := math.Sin(radians)
	return model.Vec3{X: vector.X*cosine - vector.Y*sine, Y: vector.X*sine + vector.Y*cosine, Z: vector.Z}
}

func Lerp(a, b model.Vec3, amount float64) model.Vec3 {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	return model.Vec3{X: a.X + (b.X-a.X)*amount, Y: a.Y + (b.Y-a.Y)*amount, Z: a.Z + (b.Z-a.Z)*amount}
}

func RingPosition(index, total int, radius, height float64) model.Vec3 {
	if total <= 0 {
		return model.Vec3{}
	}
	angle := (2 * math.Pi * float64(index)) / float64(total)
	return model.Vec3{X: math.Cos(angle) * radius, Y: math.Sin(angle) * radius, Z: height}
}

func Cross(a, b model.Vec3) model.Vec3 {
	return model.Vec3{X: a.Y*b.Z - a.Z*b.Y, Y: a.Z*b.X - a.X*b.Z, Z: a.X*b.Y - a.Y*b.X}
}

func Project(vector, axis model.Vec3) model.Vec3 {
	normal := Normalize(axis)
	return Scale(normal, vector.Dot(normal))
}

func ClampMagnitude(vector model.Vec3, maximum float64) model.Vec3 {
	if maximum < 0 {
		return model.Vec3{}
	}
	length := Length(vector)
	if length <= maximum {
		return vector
	}
	return Scale(vector, maximum/length)
}
