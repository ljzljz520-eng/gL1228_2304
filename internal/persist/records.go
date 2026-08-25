package persist

import (
	"encoding/json"
	"fmt"
	"stagebeam/internal/model"
)

func encode(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode record: %w", err)
	}
	return data, nil
}

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return fmt.Errorf("record is empty")
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode record: %w", err)
	}
	return nil
}

func showToRecord(show model.Show) model.ShowRecord {
	return model.ShowRecord{ShowID: show.ID, Name: show.Name, Settings: show.Settings, Gesture: show.Gesture, Active: show.Active, Frame: show.Frame, Layers: model.CloneLayers(show.Layers)}
}

func recordToShow(record model.ShowRecord) model.Show {
	return model.Show{ID: record.ShowID, Name: record.Name, Settings: record.Settings, Gesture: record.Gesture, Active: record.Active, Frame: record.Frame, Layers: model.CloneLayers(record.Layers)}
}
