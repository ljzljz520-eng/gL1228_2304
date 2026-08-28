package geometry

import (
	"math"
	"stagebeam/internal/model"
)

func GestureDirection(gesture model.Gesture, strength float64) model.Vec3 {
	if strength < 0 {
		strength = 0
	}
	if strength > 1 {
		strength = 1
	}
	switch gesture {
	case model.GestureFist:
		return model.Vec3{Z: 0.35 + strength*0.65}
	case model.GestureOpen:
		return model.Vec3{X: 0.8 * strength, Y: 0.25}
	case model.GestureRing:
		return model.Vec3{X: -0.4, Y: 0.4, Z: 0.7}
	case model.GestureWave:
		return model.Vec3{X: math.Sin(strength * math.Pi), Y: 0.65}
	default:
		return model.Vec3{Z: 0.1}
	}
}

func FanPoint(index, total int, spread, strength float64) model.Vec3 {
	if total <= 1 {
		return model.Vec3{Y: strength}
	}
	center := float64(total-1) / 2
	angle := (float64(index) - center) * spread / float64(total-1)
	return model.Vec3{X: math.Sin(angle) * strength, Y: math.Cos(angle) * strength, Z: math.Abs(math.Sin(angle))}
}

func BeamPoint(layer, total int, intensity float64) model.Vec3 {
	if total <= 0 {
		return model.Vec3{}
	}
	radius := 0.08 + float64(layer+1)/float64(total)*0.55
	return RingPosition(layer, total, radius, intensity*(1-radius))
}

func RotateParticles(particles []model.Particle, radians float64) []model.Particle {
	rotated := make([]model.Particle, len(particles))
	for index, particle := range particles {
		particle.Position = RotateZ(particle.Position, radians)
		particle.Velocity = RotateZ(particle.Velocity, radians*0.5)
		rotated[index] = particle
	}
	return rotated
}

func SpiralPosition(index, total int, radius, lift float64) model.Vec3 {
	if total <= 0 {
		return model.Vec3{}
	}
	fraction := float64(index) / float64(total)
	angle := fraction * math.Pi * 6
	return model.Vec3{X: math.Cos(angle) * radius * (0.5 + fraction*0.5), Y: math.Sin(angle) * radius * (0.5 + fraction*0.5), Z: lift * fraction}
}

func Orbit(position model.Vec3, center model.Vec3, radians float64) model.Vec3 {
	relative := model.Vec3{X: position.X - center.X, Y: position.Y - center.Y, Z: position.Z - center.Z}
	rotated := RotateZ(relative, radians)
	return Add(center, rotated)
}

func Distance(a, b model.Vec3) float64 {
	delta := model.Vec3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
	return Length(delta)
}
