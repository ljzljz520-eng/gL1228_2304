package model

import (
	"encoding/json"
	"fmt"
)

func EncodeShow(show Show) ([]byte, error) {
	if err := ValidateShow(show); err != nil {
		return nil, err
	}
	return json.Marshal(show)
}

func DecodeShow(data []byte) (Show, error) {
	var show Show
	if err := json.Unmarshal(data, &show); err != nil {
		return show, fmt.Errorf("decode show: %w", err)
	}
	if err := ValidateShow(show); err != nil {
		return show, err
	}
	return show, nil
}

func EncodeFrame(frame RenderFrame) ([]byte, error) {
	if frame.ShowID == "" {
		return nil, fmt.Errorf("frame show id is required")
	}
	return json.Marshal(frame)
}

func DecodeFrame(data []byte) (RenderFrame, error) {
	var frame RenderFrame
	if err := json.Unmarshal(data, &frame); err != nil {
		return frame, fmt.Errorf("decode frame: %w", err)
	}
	if frame.ShowID == "" {
		return frame, fmt.Errorf("frame show id is required")
	}
	return frame, nil
}

func SettingsFingerprint(settings Settings) string {
	return fmt.Sprintf("%s/%d/%d/%.3f/%.3f", settings.Color.String(), settings.ParticleCount, settings.BeamLayers, settings.Spread, settings.Intensity)
}

func LayerFingerprint(layers []BeamLayer) string {
	if len(layers) == 0 {
		return "empty"
	}
	result := ""
	for _, layer := range layers {
		result += fmt.Sprintf("%d:%.3f:%.3f:%s;", layer.Index, layer.Width, layer.Brightness, layer.Color.String())
	}
	return result
}
