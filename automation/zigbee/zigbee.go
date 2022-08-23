package zigbee

import (
	"context"
	"github.com/peter-mount/go-kernel/util/task"
	"github.com/peter-mount/home-automation/automation/model"
	"github.com/peter-mount/home-automation/automation/state"
	"github.com/peter-mount/home-automation/util/mq"
	"log"
	"strings"
	"sync"
)

// Zigbee service listens for messages from Zigbee2MQTT and passes them to the house model to see what actions
// should be performed
type Zigbee struct {
	mq        *mq.MQ         `kernel:"inject"`
	states    state.Service  `kernel:"inject"`
	house     *model.Service `kernel:"inject"`
	queueName *mq.Queue      `kernel:"config,automationQueue"`
	modelFile *string        `kernel:"config,modelFile"`
	worker    task.Queue     `kernel:"worker"`
	mutex     sync.Mutex
}

func (s *Zigbee) Start() error {
	return s.mq.ConsumeTask(s.queueName, "automation", mq.Guard(s.processZigbeeMessage))
}

// processZigbeeMessage processes a message received from zigbee2mqtt
func (s *Zigbee) processZigbeeMessage(ctx context.Context) error {
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

	s.worker.AddTask(task.Of(s.house.ScanAutomations).
		WithValue(mq.DeliveryKey, msg).
		WithValue(state.ServiceKey, s.states).
		WithValue(state.Key, newState).
		WithValue(state.PreviousStateKey, previousState))

	return nil
}
