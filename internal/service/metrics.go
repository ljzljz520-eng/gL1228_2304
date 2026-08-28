package service

import (
	"fmt"
	"stagebeam/internal/model"
)

type Metrics struct {
	ShowID        string  `json:"show_id"`
	Frame         int64   `json:"frame"`
	LayerCount    int     `json:"layer_count"`
	ParticleCount int     `json:"particle_count"`
	Energy        float64 `json:"energy"`
	Active        bool    `json:"active"`
}

func (service *ShowService) Metrics(id string) (Metrics, error) {
	show, err := service.LoadShow(id)
	if err != nil {
		return Metrics{}, err
	}
	energy := 0.0
	for _, layer := range show.Layers {
		energy += layer.Energy()
	}
	for _, particle := range show.Particles {
		energy += particle.Energy * particle.Size
	}
	return Metrics{ShowID: show.ID, Frame: show.Frame, LayerCount: len(show.Layers), ParticleCount: len(show.Particles), Energy: energy, Active: show.Active}, nil
}

func (service *ShowService) ValidateRuntime(id string) error {
	show, err := service.LoadShow(id)
	if err != nil {
		return err
	}
	if err := model.ValidateShow(show); err != nil {
		return err
	}
	if len(show.Layers) != show.Settings.BeamLayers && show.Gesture != model.GestureRing {
		return fmt.Errorf("show has %d layers, settings require %d", len(show.Layers), show.Settings.BeamLayers)
	}
	return nil
}
