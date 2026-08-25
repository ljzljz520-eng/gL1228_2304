package render

import (
	"encoding/json"
	"fmt"
	"stagebeam/internal/model"
)

type Export struct {
	Show      string             `json:"show"`
	Frame     int64              `json:"frame"`
	Message   string             `json:"message"`
	Energy    float64            `json:"energy"`
	Palette   []model.Color      `json:"palette"`
	LayerRows []LayerDescription `json:"layers"`
}

type LayerDescription struct {
	Index      int     `json:"index"`
	Width      float64 `json:"width"`
	Brightness float64 `json:"brightness"`
	Color      string  `json:"color"`
}

func ExportFrame(frame model.RenderFrame) Export {
	rows := make([]LayerDescription, 0, len(frame.Layers))
	for _, layer := range frame.Layers {
		rows = append(rows, LayerDescription{Index: layer.Index, Width: layer.Width, Brightness: layer.Brightness, Color: layer.Color.String()})
	}
	return Export{Show: frame.ShowID, Frame: frame.Frame, Message: frame.Message, Energy: frame.TotalEnergy(), Palette: FramePalette(frame), LayerRows: rows}
}

func EncodeExport(frame model.RenderFrame) ([]byte, error) {
	data, err := json.MarshalIndent(ExportFrame(frame), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode frame export: %w", err)
	}
	return data, nil
}
