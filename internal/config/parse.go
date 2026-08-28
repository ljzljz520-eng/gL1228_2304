package config

import (
	"encoding/json"
	"fmt"
	"stagebeam/internal/model"
	"strings"
)

type Request struct {
	Name     string         `json:"name"`
	Preset   string         `json:"preset"`
	Settings model.Settings `json:"settings"`
}

func ParseRequest(payload []byte) (Request, error) {
	var request Request
	if len(strings.TrimSpace(string(payload))) == 0 {
		return request, fmt.Errorf("request payload is empty")
	}
	if err := json.Unmarshal(payload, &request); err != nil {
		return request, fmt.Errorf("parse show request: %w", err)
	}
	if err := model.ValidateShowName(request.Name); err != nil {
		return request, err
	}
	if request.Preset != "" {
		preset, err := SettingsForPreset(request.Preset)
		if err != nil {
			return request, err
		}
		request.Settings = preset
	} else {
		request.Settings = ClampSettings(request.Settings)
		if err := model.ValidateSettings(request.Settings); err != nil {
			return request, err
		}
	}
	return request, nil
}

func EncodeRequest(request Request) ([]byte, error) {
	if err := model.ValidateShowName(request.Name); err != nil {
		return nil, err
	}
	if err := model.ValidateSettings(request.Settings); err != nil {
		return nil, err
	}
	return json.Marshal(request)
}

func MergeRequest(base Request, override Request) Request {
	merged := base
	if strings.TrimSpace(override.Name) != "" {
		merged.Name = override.Name
	}
	if override.Preset != "" {
		merged.Preset = override.Preset
	}
	if override.Settings != (model.Settings{}) {
		merged.Settings = override.Settings
	}
	return merged
}

func ParseSettings(payload []byte) (model.Settings, error) {
	var settings model.Settings
	if err := json.Unmarshal(payload, &settings); err != nil {
		return model.Settings{}, fmt.Errorf("parse settings: %w", err)
	}
	settings = ClampSettings(settings)
	if err := model.ValidateSettings(settings); err != nil {
		return model.Settings{}, err
	}
	return settings, nil
}
