package state

import (
	"context"
	"strings"
	"sync"
)

type Service struct {
	mutex   sync.Mutex
	states  map[string]*State // Current state of devices as reported by zigbee2mqtt
	globals map[string]string // Global state used in automation
}

func (s *Service) Start() error {
	s.states = map[string]*State{}
	s.globals = map[string]string{}
	return nil
}

func FixKey(k string) string {
	return strings.ReplaceAll(strings.ReplaceAll(k, "/", "."), " ", ".")
}

// GetState returns the latest state for a device
func (s *Service) GetState(key string) *State {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.states[FixKey(key)]
}

// SetState atomically replaces the state for a device, returning the previous one if present.
// If this returns nil then this is the first time a device has been seen
func (s *Service) SetState(key string, state *State) *State {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key = FixKey(key)
	previousState := s.states[key]
	s.states[key] = state
	return previousState
}

// GetService returns the Service instance within the context
func GetService(ctx context.Context) *Service {
	return ctx.Value(ServiceKey).(*Service)
}

func (s *Service) GetGlobal(name string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.globals[name]
}

func (s *Service) SetGlobal(name, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.globals[name] = value
}
