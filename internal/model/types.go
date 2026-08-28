package model

import (
	"fmt"
	"strings"
	"time"
)

type Gesture string

const (
	GestureOpen Gesture = "open"
	GestureFist Gesture = "fist"
	GestureRing Gesture = "ring"
	GestureWave Gesture = "wave"
	GestureIdle Gesture = "idle"
)

type Color struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

func (c Color) WithAlpha(alpha uint8) Color {
	c.A = alpha
	return c
}

func (c Color) IsDark() bool {
	return int(c.R)+int(c.G)+int(c.B) < 210
}

func (c Color) Luminance() float64 {
	return 0.2126*float64(c.R)/255 + 0.7152*float64(c.G)/255 + 0.0722*float64(c.B)/255
}

func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

func (v Vec3) Dot(other Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

type Vec3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Settings struct {
	Color         Color   `json:"color"`
	ParticleCount int     `json:"particle_count"`
	Background    Color   `json:"background"`
	BeamLayers    int     `json:"beam_layers"`
	Spread        float64 `json:"spread"`
	Intensity     float64 `json:"intensity"`
}

func (settings Settings) Label() string {
	return fmt.Sprintf("%d particles · %d layers · spread %.2f", settings.ParticleCount, settings.BeamLayers, settings.Spread)
}

func (settings Settings) Equal(other Settings) bool {
	return settings.Color == other.Color && settings.ParticleCount == other.ParticleCount && settings.Background == other.Background && settings.BeamLayers == other.BeamLayers && settings.Spread == other.Spread && settings.Intensity == other.Intensity
}

type BeamLayer struct {
	Index      int     `json:"index"`
	Width      float64 `json:"width"`
	Brightness float64 `json:"brightness"`
	Color      Color   `json:"color"`
	Angle      float64 `json:"angle"`
}

func (layer BeamLayer) Energy() float64 {
	return layer.Width * layer.Brightness
}

func (layer BeamLayer) WithColor(color Color) BeamLayer {
	layer.Color = color
	return layer
}

type Particle struct {
	ID       int     `json:"id"`
	Position Vec3    `json:"position"`
	Velocity Vec3    `json:"velocity"`
	Energy   float64 `json:"energy"`
	Size     float64 `json:"size"`
}

type Show struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Settings  Settings    `json:"settings"`
	Gesture   Gesture     `json:"gesture"`
	Active    bool        `json:"active"`
	Frame     int64       `json:"frame"`
	Layers    []BeamLayer `json:"layers"`
	Particles []Particle  `json:"particles"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type GestureEvent struct {
	ShowID   string  `json:"show_id"`
	Gesture  Gesture `json:"gesture"`
	Strength float64 `json:"strength"`
	Ended    bool    `json:"ended"`
	Sequence int64   `json:"sequence"`
}

func (event GestureEvent) Label() string {
	if event.Ended {
		return string(event.Gesture) + " ended"
	}
	return string(event.Gesture) + " active"
}

type RenderFrame struct {
	ShowID     string      `json:"show_id"`
	Frame      int64       `json:"frame"`
	Gesture    Gesture     `json:"gesture"`
	Background Color       `json:"background"`
	Layers     []BeamLayer `json:"layers"`
	Particles  []Particle  `json:"particles"`
	Message    string      `json:"message"`
}

func (frame RenderFrame) TotalEnergy() float64 {
	energy := 0.0
	for _, layer := range frame.Layers {
		energy += layer.Energy()
	}
	for _, particle := range frame.Particles {
		energy += particle.Energy * particle.Size
	}
	return energy
}

func (frame RenderFrame) Summary() string {
	return fmt.Sprintf("%s: frame %d, %d layers, %d particles", frame.ShowID, frame.Frame, len(frame.Layers), len(frame.Particles))
}

type AuditEntry struct {
	ID        string    `json:"id"`
	ShowID    string    `json:"show_id"`
	Action    string    `json:"action"`
	Detail    string    `json:"detail"`
	Sequence  int64     `json:"sequence"`
	CreatedAt time.Time `json:"created_at"`
}

type ShowRecord struct {
	ShowID   string      `json:"show_id"`
	Name     string      `json:"name"`
	Settings Settings    `json:"settings"`
	Gesture  Gesture     `json:"gesture"`
	Active   bool        `json:"active"`
	Frame    int64       `json:"frame"`
	Layers   []BeamLayer `json:"layers"`
}

type GestureRecord struct {
	ShowID   string  `json:"show_id"`
	Gesture  Gesture `json:"gesture"`
	Strength float64 `json:"strength"`
	Ended    bool    `json:"ended"`
	Sequence int64   `json:"sequence"`
}

type LayerSnapshot struct {
	ShowID string      `json:"show_id"`
	Frame  int64       `json:"frame"`
	Layers []BeamLayer `json:"layers"`
}

type SettingRecord struct {
	ShowID   string   `json:"show_id"`
	Settings Settings `json:"settings"`
}

type WorkflowReceipt struct {
	ShowID     string `json:"show_id"`
	Operation  string `json:"operation"`
	Frame      int64  `json:"frame"`
	LayerCount int    `json:"layer_count"`
	AuditID    string `json:"audit_id"`
}

func (receipt WorkflowReceipt) Complete() bool {
	return strings.TrimSpace(receipt.ShowID) != "" && receipt.Operation != "" && receipt.AuditID != ""
}

func CloneShow(show Show) Show {
	show.Layers = CloneLayers(show.Layers)
	show.Particles = CloneParticles(show.Particles)
	return show
}

func (show Show) DisplayName() string {
	if strings.TrimSpace(show.Name) == "" {
		return show.ID
	}
	return strings.TrimSpace(show.Name)
}
