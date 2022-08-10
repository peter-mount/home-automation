package cache

import (
	"strings"
	"time"
)

// State message for common settings.
//
// See https://www.zigbee2mqtt.io/guide/usage/mqtt_topics_and_messages.html
//
type State struct {
	StateUpdate
	LastSeen           *time.Time             `json:"last_seen,omitempty"`           // timestamp
	Battery            uint8                  `json:"battery,omitempty"`             // Battery level
	BatteryLow         bool                   `json:"battery_low,omitempty"`         // battery low warning
	Voltage            int                    `json:"voltage,omitempty"`             // Battery voltage
	AirQuality         string                 `json:"air_quality,omitempty"`         // Air quality
	Voc                int                    `json:"voc,omitempty"`                 // Air quality
	Temperature        float64                `json:"temperature,omitempty"`         // temp C
	Humidity           float64                `json:"humidity,omitempty"`            // Humidity %
	LinkQuality        uint8                  `json:"link_quality,omitempty"`        // Link quality
	ColourMode         string                 `json:"color_mode,omitempty"`          // Colour mode
	Colour             *Colour                `json:"color,omitempty"`               // Colour
	Smoke              bool                   `json:"smoke,omitempty"`               // fire alarm detection
	RestoreReports     bool                   `json:"restore_reports,omitempty"`     // fire alarm
	SupervisionReports bool                   `json:"supervision_reports,omitempty"` // fire alarm
	Test               bool                   `json:"test,omitempty"`                // Test triggered
	Action             string                 `json:"action,omitempty"`              // Last Button pressed
	Contact            bool                   `json:"contact,omitempty"`             // Window/Door is closed (false = open)
	Tamper             bool                   `json:"tamper,omitempty"`              // Sensor, tamper triggered
	Update             map[string]interface{} `json:"update,omitempty"`              // Misc data
	UpdateAvailable    bool                   `json:"update_available,omitempty"`    // Firmware update available
}

// StateUpdate contains fields writable to zigbee2mqtt
type StateUpdate struct {
	parent     *State // used for linking to old data, see toggle
	State      string `json:"state,omitempty"`      // State, either ON or OFF
	Brightness uint8  `json:"brightness,omitempty"` // Brightness, 0..255
	ColourTemp uint   `json:"color_temp,omitempty"` // Colour Temperature
}

type Colour struct {
	X          float64 `json:"x,omitempty"`          // Colour X component
	Y          float64 `json:"y,omitempty"`          // Colour Y component
	Hue        float64 `json:"hue,omitempty"`        // Hue
	Saturation float64 `json:"saturation,omitempty"` // Saturation
}

func (s *State) DoUpdate() *StateUpdate {
	return &StateUpdate{parent: s}
}

func (s *StateUpdate) ensureUpdatable() {
	if s.parent == nil {
		panic("must use Update() first!")
	}
}

// Toggle state between ON and OFF
func (s *StateUpdate) Toggle() *StateUpdate {
	s.ensureUpdatable()
	// Check original state
	if strings.ToUpper(s.parent.State) == "ON" {
		s.State = "OFF"
	} else {
		s.State = "ON"
	}
	return s
}
func (s *StateUpdate) SetState(state string) *StateUpdate {
	s.ensureUpdatable()
	s.State = strings.ToUpper(state)
	if s.State == "TOGGLE" {
		return s.Toggle()
	} else if s.State == "" {
		s.State = "ON"
	}
	return s
}
