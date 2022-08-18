package util

import (
	cache2 "github.com/peter-mount/home-automation/automation/cache"
	"github.com/peter-mount/home-automation/util/mq"
	"time"
)

// Event issued to RabbitMQ with data on what's happened to a device.
// This is also used as the response to get requests on a device
type Event struct {
	Name      string         `json:"name,omitempty"`
	Routing   string         `json:"routing,omitempty"`
	Type      string         `json:"type,omitempty"`
	Timestamp time.Time      `json:"timestamp,omitempty"`
	State     *cache2.State  `json:"state,omitempty"`
	Device    *cache2.Device `json:"device,omitempty"`
}

func NewEvent(name, eventType string, state *cache2.State, device *cache2.Device) *Event {
	var ts time.Time

	if state != nil && state.LastSeen != nil && !state.LastSeen.IsZero() {
		ts = *state.LastSeen
	}

	if device != nil {
		if name == "" {
			name = device.FriendlyName
		}
		if device.LastSeen > 0 {
			t2 := time.Unix(device.LastSeen/1000, 0)
			if t2.After(ts) {
				ts = t2
			}
		}
	}

	if ts.IsZero() {
		ts = time.Now()
	}

	return &Event{
		Name:      name,
		Routing:   mq.EncodeKey(name),
		Timestamp: ts,
		Type:      eventType,
		State:     state,
		Device:    device,
	}
}
