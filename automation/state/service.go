package state

import (
	"github.com/peter-mount/go-kernel"
	"sync"
)

type Service interface {
	// GetState returns the latest state for a device
	GetState(key string) *State

	// SetState atomically replaces the state for a device, returning the previous one if present.
	// If this returns nil then this is the first time a device has been seen
	SetState(key string, state *State) *State

	// GetGlobal returns the current value of a global variable
	GetGlobal(name string) string

	// SetGlobal sets a global variable
	SetGlobal(name, value string)
}

func init() {
	kernel.RegisterAPI((*Service)(nil), &service{})
}

type service struct {
	mutex   sync.Mutex
	states  map[string]*State // Current state of devices as reported by zigbee2mqtt
	globals map[string]string // Global state used in automation
}

func (s *service) Start() error {
	s.states = map[string]*State{}
	s.globals = map[string]string{}
	return nil
}

// GetState returns the latest state for a device
func (s *service) GetState(key string) *State {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.states[FixKey(key)]
}

// SetState atomically replaces the state for a device, returning the previous one if present.
// If this returns nil then this is the first time a device has been seen
func (s *service) SetState(key string, newState *State) *State {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key = FixKey(key)
	previousState := s.states[key]
	s.states[key] = newState
	return previousState
}

// GetGlobal returns the current value of a global variable
func (s *service) GetGlobal(name string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.globals[name]
}

// SetGlobal sets a global variable
func (s *service) SetGlobal(name, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.globals[name] = value
}
