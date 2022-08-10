package automation

import (
	"context"
	"github.com/peter-mount/go-kernel/util/task"
	"github.com/peter-mount/home-automation/model"
	"github.com/peter-mount/home-automation/state"
	"log"
)

func (s *Service) scanAutomations(ctx context.Context) error {
	for autoID, auto := range s.GetModel().Automation {
		if auto.IsTriggered(ctx) {
			log.Printf("Triggered Automation %q", autoID)
			s.worker.AddTask(task.Of(s.triggerAutomation).WithValue("automation", auto))
		}
	}
	return nil
}

func (s *Service) triggerAutomation(ctx context.Context) error {
	auto := ctx.Value("automation").(*model.Automation)
	for _, action := range auto.Actions {

		// Activate named scene if it exists
		if action.Scene != "" {
			if scene, exists := s.GetModel().Scene[action.Scene]; exists {
				s.worker.AddTask(task.Of(s.activateScene).WithValue("scene", scene))
			}
		}

		// Set any global variables
		if action.Global != nil {
			for k, v := range *action.Global {
				s.states.SetGlobal(k, v)
			}
		}
	}
	return nil
}

func (s *Service) activateScene(ctx context.Context) error {
	house := s.GetModel()
	scene := ctx.Value("scene").(*model.Scene)
	for deviceID, stateId := range scene.Devices {
		if device, exists := house.Device[deviceID]; exists {
			if currentState, exists := device.Action[stateId]; exists {
				s.worker.AddTask(task.Of(s.activateDevice).
					WithValue("deviceID", deviceID).
					WithValue("state", currentState))
			}
		}
	}
	return nil
}

func (s *Service) activateDevice(ctx context.Context) error {
	deviceID := ctx.Value("deviceID").(string)
	for _, currentState := range ctx.Value("state").([]*state.State) {
		err := s.publisher.PublishJSON(deviceID+".set", currentState)
		if err != nil {
			return err
		}
	}
	return nil
}
