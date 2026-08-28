package transport

import (
	"encoding/json"
	"fmt"
	"stagebeam/internal/config"
	"stagebeam/internal/model"
	"stagebeam/internal/service"
)

type Command struct {
	Action   string          `json:"action"`
	ShowID   string          `json:"show_id"`
	Name     string          `json:"name"`
	Preset   string          `json:"preset"`
	Gesture  model.Gesture   `json:"gesture"`
	Strength float64         `json:"strength"`
	Ended    bool            `json:"ended"`
	Settings *model.Settings `json:"settings"`
}

type Response struct {
	OK      bool                   `json:"ok"`
	Message string                 `json:"message"`
	Frame   *model.RenderFrame     `json:"frame,omitempty"`
	Receipt *model.WorkflowReceipt `json:"receipt,omitempty"`
}

func DecodeCommand(data []byte) (Command, error) {
	var command Command
	if err := json.Unmarshal(data, &command); err != nil {
		return command, fmt.Errorf("decode command: %w", err)
	}
	if command.Action == "" {
		return command, fmt.Errorf("command action is required")
	}
	if command.Strength == 0 && command.Gesture != model.GestureIdle {
		command.Strength = 0.5
	}
	return command, nil
}

func Execute(service *service.ShowService, command Command) Response {
	switch command.Action {
	case "create":
		settings := config.DefaultSettings()
		if command.Preset != "" {
			preset, err := config.SettingsForPreset(command.Preset)
			if err != nil {
				return failure(err)
			}
			settings = preset
		}
		receipt, err := service.CreateFromRequest(command.ShowID, config.Request{Name: command.Name, Settings: settings, Preset: command.Preset})
		if err != nil {
			return failure(err)
		}
		return Response{OK: true, Message: "show created", Receipt: &receipt}
	case "gesture":
		frame, receipt, err := service.ApplyGesture(model.GestureEvent{ShowID: command.ShowID, Gesture: command.Gesture, Strength: command.Strength, Ended: command.Ended})
		if err != nil {
			return failure(err)
		}
		return Response{OK: true, Message: frame.Message, Frame: &frame, Receipt: &receipt}
	case "frame":
		frame, err := service.Frame(command.ShowID)
		if err != nil {
			return failure(err)
		}
		return Response{OK: true, Message: frame.Message, Frame: &frame}
	case "close":
		receipt, err := service.CloseShow(command.ShowID)
		if err != nil {
			return failure(err)
		}
		return Response{OK: true, Message: "show closed", Receipt: &receipt}
	default:
		return failure(fmt.Errorf("unknown action %q", command.Action))
	}
}

func EncodeResponse(response Response) ([]byte, error) {
	return json.Marshal(response)
}

func failure(err error) Response {
	return Response{OK: false, Message: err.Error()}
}

func CommandForCreate(id, name, preset string) Command {
	return Command{Action: "create", ShowID: id, Name: name, Preset: preset}
}

func CommandForGesture(id string, gesture model.Gesture, strength float64, ended bool) Command {
	return Command{Action: "gesture", ShowID: id, Gesture: gesture, Strength: strength, Ended: ended}
}

func CommandForFrame(id string) Command {
	return Command{Action: "frame", ShowID: id}
}

func ValidateCommand(command Command) error {
	if command.ShowID == "" && command.Action != "create" {
		return fmt.Errorf("show id is required for %s", command.Action)
	}
	if command.Action == "create" && command.Name == "" {
		return fmt.Errorf("name is required for create")
	}
	if command.Action == "gesture" && (command.Strength < 0 || command.Strength > 1) {
		return fmt.Errorf("gesture strength must be between zero and one")
	}
	return nil
}

func ExecuteJSON(showService *service.ShowService, data []byte) ([]byte, error) {
	command, err := DecodeCommand(data)
	if err != nil {
		return nil, err
	}
	if err := ValidateCommand(command); err != nil {
		return nil, err
	}
	response := Execute(showService, command)
	return EncodeResponse(response)
}
