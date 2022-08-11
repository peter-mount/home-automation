package model

import (
	"context"
	"github.com/peter-mount/go-kernel/util/task"
	"github.com/peter-mount/home-automation/state"
	automation "github.com/peter-mount/home-automation/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// Service is a Kernel service for managing the House model
type Service struct {
	modelFile *string               `kernel:"config,modelFile"`
	states    *state.Service        `kernel:"inject"`
	publisher *automation.Publisher `kernel:"inject"`
	worker    task.Queue            `kernel:"worker"`
	mutex     sync.Mutex
	house     *House
}

func (s *Service) Start() error {
	return s.LoadModel()
}

func (s *Service) LoadModel() error {
	f, err := os.Open(*s.modelFile)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	house := &House{}
	err = yaml.Unmarshal(b, house)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.house = house
	return nil
}

// GetModel gets the current House instance.
// Never access house direct, but always via this function as it could be reloaded at any time.
// It is not safe to cache the instance returned unless it's within the scope of a function.
func (s *Service) GetModel() *House {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.house
}

// ScanAutomations scans all automations within the House for triggered events
func (s *Service) ScanAutomations(ctx context.Context) error {
	for autoID, auto := range s.GetModel().Automation {
		if auto.IsTriggered(ctx) {
			log.Printf("Triggered Automation %q", autoID)
			s.worker.AddTask(task.Of(s.triggerAutomation).WithValue("automation", auto))
		}
	}
	return nil
}

// triggerAutomation is invoked when an Automation has been triggered by an incoming event.
// Context Keys: "automation" with the *Automation that has been triggered
func (s *Service) triggerAutomation(ctx context.Context) error {
	auto := ctx.Value("automation").(*Automation)

	for _, action := range auto.Actions {
		s.runAction(action)
	}

	return nil
}

// runAction is called from triggerAutomation for each action
func (s *Service) runAction(action *Action) {
	house := s.GetModel()

	// Activate named scene if it exists
	if action.Scene != "" {
		if scene, exists := house.Scene[action.Scene]; exists {

			for deviceID, stateId := range scene.Devices {
				if device, exists := house.Device[deviceID]; exists {
					if currentState, exists := device.Action[stateId]; exists {
						s.activateDevice(deviceID, currentState)
					}
				}
			}
		}
	}

	// Set any global variables
	if action.Global != nil {
		for k, v := range *action.Global {
			s.states.SetGlobal(k, v)
		}
	}
}

func (s *Service) activateDevice(deviceID string, states []*state.State) {
	defer func() {
		_ = recover()
	}()

	for _, currentState := range states {
		destination := deviceID + ".set"
		err := s.publisher.PublishJSON(destination, currentState)
		if err != nil {
			log.Printf("Failed to publish to %s: %v", destination, err)
		}
	}
}
