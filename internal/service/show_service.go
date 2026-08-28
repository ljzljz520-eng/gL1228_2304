package service

import (
	"fmt"
	"stagebeam/internal/audit"
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"stagebeam/internal/persist"
	"stagebeam/internal/render"
	"stagebeam/internal/stage"
	"sync"
)

type ShowService struct {
	store    *persist.Store
	audit    *audit.Log
	composer stage.LayerComposer
	renderer render.Renderer
	mu       sync.Mutex
	shows    map[string]model.Show
}

func NewShowService(store *persist.Store) *ShowService {
	return &ShowService{store: store, audit: audit.NewLog(), composer: stage.NewLayerComposer(), renderer: render.NewRenderer(512), shows: make(map[string]model.Show)}
}

func (service *ShowService) CreateShow(id, name string, settings model.Settings) (model.Show, error) {
	service.mu.Lock()
	defer service.mu.Unlock()
	if id == "" {
		return model.Show{}, fmt.Errorf("show id is required")
	}
	if err := model.ValidateShowName(name); err != nil {
		return model.Show{}, err
	}
	settings = config.ClampSettings(settings)
	if err := model.ValidateSettings(settings); err != nil {
		return model.Show{}, err
	}
	show := model.Show{ID: id, Name: name, Settings: settings, Gesture: model.GestureIdle, Active: true, Layers: service.composer.Build(settings, model.GestureIdle, 0), Particles: stage.SeedParticles(settings)}
	service.shows[id] = show
	if service.store != nil {
		if err := service.store.SaveShow(show); err != nil {
			delete(service.shows, id)
			return model.Show{}, err
		}
	}
	service.audit.Record(id, "show.created", name)
	return cloneShow(show), nil
}

func (service *ShowService) LoadShow(id string) (model.Show, error) {
	service.mu.Lock()
	defer service.mu.Unlock()
	if show, ok := service.shows[id]; ok {
		return cloneShow(show), nil
	}
	if service.store == nil {
		return model.Show{}, persist.ErrNotFound
	}
	show, err := service.store.LoadShow(id)
	if err != nil {
		return model.Show{}, err
	}
	show.Particles = stage.SeedParticles(show.Settings)
	service.shows[id] = show
	return cloneShow(show), nil
}

func (service *ShowService) UpdateSettings(id string, settings model.Settings) (model.WorkflowReceipt, error) {
	service.mu.Lock()
	defer service.mu.Unlock()
	show, ok := service.shows[id]
	if !ok {
		if service.store == nil {
			return model.WorkflowReceipt{}, persist.ErrNotFound
		}
		var err error
		show, err = service.store.LoadShow(id)
		if err != nil {
			return model.WorkflowReceipt{}, err
		}
	}
	settings = config.ClampSettings(settings)
	if err := model.ValidateSettings(settings); err != nil {
		return model.WorkflowReceipt{}, err
	}
	show.Settings = settings
	show.Layers = service.composer.Build(settings, show.Gesture, 0.5)
	show.Particles = stage.RecolorParticles(show.Particles, settings.Color)
	service.shows[id] = show
	if service.store != nil {
		if err := service.store.SaveShow(show); err != nil {
			return model.WorkflowReceipt{}, err
		}
	}
	entry := service.audit.Record(id, "settings.updated", fmt.Sprintf("layers=%d particles=%d", settings.BeamLayers, settings.ParticleCount))
	return model.WorkflowReceipt{ShowID: id, Operation: "settings.updated", Frame: show.Frame, LayerCount: len(show.Layers), AuditID: entry.ID}, nil
}

func (service *ShowService) ApplyGesture(event model.GestureEvent) (model.RenderFrame, model.WorkflowReceipt, error) {
	service.mu.Lock()
	defer service.mu.Unlock()
	if err := model.ValidateEvent(event); err != nil {
		return model.RenderFrame{}, model.WorkflowReceipt{}, err
	}
	show, ok := service.shows[event.ShowID]
	if !ok {
		if service.store == nil {
			return model.RenderFrame{}, model.WorkflowReceipt{}, persist.ErrNotFound
		}
		var err error
		show, err = service.store.LoadShow(event.ShowID)
		if err != nil {
			return model.RenderFrame{}, model.WorkflowReceipt{}, err
		}
		show.Particles = stage.SeedParticles(show.Settings)
	}
	gesture := model.NormalizeGesture(event.Gesture)
	show.Gesture = gesture
	show.Frame++
	if gesture == model.GestureOpen {
		show.Layers = stage.FanLayers(show.Settings, event.Strength)
	} else {
		show.Layers = service.composer.Build(show.Settings, gesture, event.Strength)
	}
	show.Layers = service.composer.ComposeLayers(show.Layers, gesture, event.Ended)
	show.Particles = stage.AdvanceParticles(show.Particles, gesture, event.Strength, show.Frame)
	if event.Ended {
		show.Gesture = model.GestureIdle
	}
	service.shows[event.ShowID] = show
	if service.store != nil {
		if err := service.store.SaveShow(show); err != nil {
			return model.RenderFrame{}, model.WorkflowReceipt{}, err
		}
		if err := service.store.SaveGesture(model.GestureRecord{ShowID: event.ShowID, Gesture: gesture, Strength: event.Strength, Ended: event.Ended, Sequence: event.Sequence}); err != nil {
			return model.RenderFrame{}, model.WorkflowReceipt{}, err
		}
		if err := service.store.SaveLayerSnapshot(model.LayerSnapshot{ShowID: event.ShowID, Frame: show.Frame, Layers: show.Layers}); err != nil {
			return model.RenderFrame{}, model.WorkflowReceipt{}, err
		}
	}
	entry := service.audit.Record(event.ShowID, "gesture.applied", fmt.Sprintf("gesture=%s ended=%t", gesture, event.Ended))
	frame := service.renderer.Frame(show)
	return frame, model.WorkflowReceipt{ShowID: event.ShowID, Operation: "gesture.applied", Frame: show.Frame, LayerCount: len(show.Layers), AuditID: entry.ID}, nil
}

func (service *ShowService) CloseShow(id string) (model.WorkflowReceipt, error) {
	service.mu.Lock()
	defer service.mu.Unlock()
	show, ok := service.shows[id]
	if !ok {
		return model.WorkflowReceipt{}, persist.ErrNotFound
	}
	show.Active = false
	service.shows[id] = show
	if service.store != nil {
		if err := service.store.SaveShow(show); err != nil {
			return model.WorkflowReceipt{}, err
		}
	}
	entry := service.audit.Record(id, "show.closed", show.Name)
	return model.WorkflowReceipt{ShowID: id, Operation: "show.closed", Frame: show.Frame, LayerCount: len(show.Layers), AuditID: entry.ID}, nil
}

func (service *ShowService) Frame(id string) (model.RenderFrame, error) {
	show, err := service.LoadShow(id)
	if err != nil {
		return model.RenderFrame{}, err
	}
	return service.renderer.Frame(show), nil
}

func (service *ShowService) AuditEntries(id string) []model.AuditEntry {
	return service.audit.ForShow(id)
}

func (service *ShowService) SaveAudit(id string, action, detail string) error {
	entry := service.audit.Record(id, action, detail)
	if service.store == nil {
		return nil
	}
	return service.store.SaveAudit(entry)
}

func (service *ShowService) ListShows() ([]model.Show, error) {
	service.mu.Lock()
	defer service.mu.Unlock()
	if service.store == nil {
		shows := make([]model.Show, 0, len(service.shows))
		for _, show := range service.shows {
			shows = append(shows, cloneShow(show))
		}
		return shows, nil
	}
	return service.store.ListShows()
}

func (service *ShowService) LayerSnapshot(id string) ([]model.BeamLayer, error) {
	service.mu.Lock()
	defer service.mu.Unlock()
	show, ok := service.shows[id]
	if !ok {
		if service.store == nil {
			return nil, persist.ErrNotFound
		}
		var err error
		show, err = service.store.LoadShow(id)
		if err != nil {
			return nil, err
		}
	}
	return model.CloneLayers(show.Layers), nil
}

func (service *ShowService) ApplyMouse(id string, point render.MousePoint, width, height float64, pressed bool) (model.RenderFrame, model.WorkflowReceipt, error) {
	gesture, strength := render.GestureFromMouse(point, width, height, pressed)
	return service.ApplyGesture(model.GestureEvent{ShowID: id, Gesture: gesture, Strength: strength, Ended: !pressed})
}

func (service *ShowService) Advance(id string, frames int64) (model.RenderFrame, error) {
	if frames < 1 {
		return service.Frame(id)
	}
	var result model.RenderFrame
	for index := int64(0); index < frames; index++ {
		frame, _, err := service.ApplyGesture(model.GestureEvent{ShowID: id, Gesture: model.GestureIdle, Strength: 0.2, Sequence: index + 1})
		if err != nil {
			return model.RenderFrame{}, err
		}
		result = frame
	}
	return result, nil
}

func cloneShow(show model.Show) model.Show {
	show.Layers = model.CloneLayers(show.Layers)
	show.Particles = model.CloneParticles(show.Particles)
	return show
}
