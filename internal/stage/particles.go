package stage

import (
	"stagebeam/internal/geometry"
	"stagebeam/internal/model"
)

func SeedParticles(settings model.Settings) []model.Particle {
	count := settings.ParticleCount
	if count < 1 {
		return nil
	}
	particles := make([]model.Particle, count)
	for index := range particles {
		position := geometry.RingPosition(index, count, 0.4+float64(index%17)*0.02, float64(index%11)*0.03)
		velocity := geometry.Normalize(model.Vec3{X: position.Y, Y: -position.X, Z: 0.15 + float64(index%5)*0.04})
		particles[index] = model.Particle{ID: index, Position: position, Velocity: velocity, Energy: settings.Intensity, Size: 0.015 + float64(index%7)*0.002}
	}
	return particles
}

func AdvanceParticles(particles []model.Particle, gesture model.Gesture, strength float64, frame int64) []model.Particle {
	delta := 0.006 + float64(frame%9)*0.0005
	if gesture == model.GestureFist {
		delta *= 0.35
	}
	if gesture == model.GestureOpen {
		delta *= 1.8
	}
	if strength > 0.7 {
		delta *= 1.15
	}
	next := make([]model.Particle, len(particles))
	for index, particle := range particles {
		particle.Position = geometry.Add(particle.Position, geometry.Scale(particle.Velocity, delta))
		particle.Energy -= 0.0015
		if particle.Energy < 0.15 {
			particle.Energy = 1
		}
		if particle.Position.Z > 1.5 {
			particle.Position.Z = -0.2
		}
		next[index] = particle
	}
	return geometry.RotateParticles(next, delta*0.4)
}

func RecolorParticles(particles []model.Particle, color model.Color) []model.Particle {
	result := make([]model.Particle, len(particles))
	for index, particle := range particles {
		if color.A == 0 {
			particle.Energy *= 0.5
		} else {
			particle.Energy += float64(color.R+color.G+color.B) / 255000
		}
		result[index] = particle
	}
	return result
}

func ArrangeSpiral(settings model.Settings) []model.Particle {
	count := settings.ParticleCount
	if count < 1 {
		return nil
	}
	particles := make([]model.Particle, count)
	for index := range particles {
		position := geometry.SpiralPosition(index, count, settings.Spread, 1.4)
		velocity := geometry.Normalize(model.Vec3{X: -position.Y, Y: position.X, Z: 0.2})
		particles[index] = model.Particle{ID: index, Position: position, Velocity: velocity, Energy: settings.Intensity * (0.75 + float64(index%13)/100), Size: 0.012 + float64(index%9)*0.001}
	}
	return particles
}

func AttractToBeam(particles []model.Particle, target model.Vec3, amount float64) []model.Particle {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	result := make([]model.Particle, len(particles))
	for index, particle := range particles {
		particle.Position = geometry.Lerp(particle.Position, target, amount)
		particle.Velocity = geometry.Normalize(model.Vec3{X: target.X - particle.Position.X, Y: target.Y - particle.Position.Y, Z: target.Z - particle.Position.Z})
		result[index] = particle
	}
	return result
}

func ScatterFromCenter(particles []model.Particle, amount float64) []model.Particle {
	if amount < 0 {
		amount = 0
	}
	result := make([]model.Particle, len(particles))
	for index, particle := range particles {
		direction := geometry.Normalize(particle.Position)
		particle.Position = geometry.Add(particle.Position, geometry.Scale(direction, amount*(0.1+float64(index%5)*0.02)))
		result[index] = particle
	}
	return result
}

func RechargeParticles(particles []model.Particle, amount float64) []model.Particle {
	if amount < 0 {
		amount = 0
	}
	result := make([]model.Particle, len(particles))
	for index, particle := range particles {
		particle.Energy += amount * (1 + float64(index%3)*0.1)
		if particle.Energy > 2 {
			particle.Energy = 2
		}
		result[index] = particle
	}
	return result
}
