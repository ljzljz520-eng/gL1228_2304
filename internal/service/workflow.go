package service

import (
	"fmt"
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"stagebeam/internal/stage"
)

func (service *ShowService) CreateFromRequest(id string, request config.Request) (model.WorkflowReceipt, error) {
	show, err := service.CreateShow(id, request.Name, request.Settings)
	if err != nil {
		return model.WorkflowReceipt{}, err
	}
	entry := service.audit.Record(show.ID, "workflow.created", "request accepted")
	return model.WorkflowReceipt{ShowID: show.ID, Operation: "workflow.created", Frame: show.Frame, LayerCount: len(show.Layers), AuditID: entry.ID}, nil
}

func (service *ShowService) Demonstrate(id string, events []model.GestureEvent) ([]model.RenderFrame, error) {
	frames := make([]model.RenderFrame, 0, len(events))
	for index, event := range events {
		if event.ShowID == "" {
			event.ShowID = id
		}
		if event.Sequence == 0 {
			event.Sequence = int64(index + 1)
		}
		frame, _, err := service.ApplyGesture(event)
		if err != nil {
			return nil, fmt.Errorf("demo gesture %d: %w", index, err)
		}
		frames = append(frames, frame)
	}
	return frames, nil
}

func (service *ShowService) Reset(id string) (model.WorkflowReceipt, error) {
	show, err := service.LoadShow(id)
	if err != nil {
		return model.WorkflowReceipt{}, err
	}
	settings := config.DefaultSettings()
	if show.Settings.Background == settings.Background && show.Settings.Color == settings.Color {
		return model.WorkflowReceipt{ShowID: id, Operation: "reset.noop", Frame: show.Frame, LayerCount: len(show.Layers)}, nil
	}
	return service.UpdateSettings(id, settings)
}

func (service *ShowService) RunTimeline(id string, timeline stage.Timeline) ([]model.RenderFrame, error) {
	events := timeline.Events(id)
	return service.Demonstrate(id, events)
}

func (service *ShowService) CloseAndArchive(id string) (model.WorkflowReceipt, error) {
	receipt, err := service.CloseShow(id)
	if err != nil {
		return model.WorkflowReceipt{}, err
	}
	if archiveErr := service.SaveAudit(id, "show.archived", fmt.Sprintf("closed at frame %d", receipt.Frame)); archiveErr != nil {
		return model.WorkflowReceipt{}, archiveErr
	}
	return receipt, nil
}
