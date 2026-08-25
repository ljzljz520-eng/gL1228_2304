package service

import (
	"fmt"
	"stagebeam/internal/model"
)

type ReplayResult struct {
	ShowID      string
	Frames      []model.RenderFrame
	FinalFrame  int64
	LayerCounts []int
}

func (service *ShowService) Replay(id string, events []model.GestureEvent) (ReplayResult, error) {
	if _, err := service.LoadShow(id); err != nil {
		return ReplayResult{}, err
	}
	result := ReplayResult{ShowID: id, Frames: make([]model.RenderFrame, 0, len(events)), LayerCounts: make([]int, 0, len(events))}
	for index, event := range events {
		event.ShowID = id
		if event.Sequence == 0 {
			event.Sequence = int64(index + 1)
		}
		frame, _, err := service.ApplyGesture(event)
		if err != nil {
			return ReplayResult{}, fmt.Errorf("replay event %d: %w", index, err)
		}
		result.Frames = append(result.Frames, frame)
		result.LayerCounts = append(result.LayerCounts, len(frame.Layers))
		result.FinalFrame = frame.Frame
	}
	return result, nil
}

func (result ReplayResult) Stable() bool {
	if len(result.Frames) == 0 {
		return false
	}
	if len(result.Frames) != len(result.LayerCounts) {
		return false
	}
	for index, frame := range result.Frames {
		if frame.Frame <= 0 || result.LayerCounts[index] != len(frame.Layers) {
			return false
		}
	}
	return true
}

func (service *ShowService) Health(id string) (Metrics, error) {
	if service.store != nil {
		if err := service.store.Health(); err != nil {
			return Metrics{}, err
		}
	}
	metrics, err := service.Metrics(id)
	if err != nil {
		return Metrics{}, err
	}
	if metrics.LayerCount == 0 || metrics.ParticleCount == 0 {
		return Metrics{}, fmt.Errorf("show %s has no renderable content", id)
	}
	return metrics, nil
}
