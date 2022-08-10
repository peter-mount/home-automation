package automation

import (
	"context"
	"github.com/peter-mount/go-kernel/util/task"
	"github.com/peter-mount/home-automation/mq"
	"github.com/peter-mount/home-automation/state"
	"log"
	"strings"
)

func (s *Service) processMqMessage(ctx context.Context) error {
	/*// Queue inbound message, place at a priority so actions take precedence
	    s.worker.AddPriorityTask(500, task.Of(s.processImpl).WithContext(ctx, mq.DeliveryKey))
	    return nil
	  }

	  func (s *Service) processImpl(ctx context.Context) error {*/
	msg := mq.Delivery(ctx)
	key := msg.RoutingKey

	// ignore bridge logging
	if strings.HasPrefix(key, "zigbee2mqtt.bridge.") {
		return nil
	}

	// HomeAssistant creating these?
	if strings.HasSuffix(key, ".action") {
		return nil
	}
	// Ignore MQTT commands
	if strings.HasSuffix(key, ".set") || strings.HasSuffix(key, ".get") {
		return nil
	}

	newState, err := state.UnmarshalState(msg.Body)
	if err != nil {
		return err
	}

	log.Printf("%q %s", msg.RoutingKey, msg.Body)

	previousState := s.states.SetState(key, newState)
	if previousState == nil {
		log.Printf("New device %q", key)
		previousState = newState
	}

	s.worker.AddTask(task.Of(s.scanAutomations).
		WithValue(mq.DeliveryKey, msg).
		WithValue(state.ServiceKey, s.states).
		WithValue(state.StateKey, newState).
		WithValue(state.PreviousStateKey, previousState))

	return nil
}
