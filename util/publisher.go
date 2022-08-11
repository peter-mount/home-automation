package util

import (
	"github.com/peter-mount/home-automation/mq"
	"log"
	"strings"
)

// Publisher handles the publishing of messages back to RabbitMQ
type Publisher struct {
	mq              *mq.MQ        `kernel:"inject"`
	zigbeePublisher *mq.Publisher `kernel:"config,automationPublisher"`
}

func (s *Publisher) Start() error {
	err := s.mq.AttachPublisher(s.zigbeePublisher)
	if err != nil {
		return err
	}

	return nil
}

func (s *Publisher) PublishJSON(key string, payload interface{}) error {
	var publisher *mq.Publisher

	// Resolve which publisher to use
	if strings.HasPrefix(key, "zigbee2mqtt/") {
		publisher = s.zigbeePublisher
	}

	if publisher != nil {
		return publisher.PublishJSON(key, payload)
	}

	// Should not happen unless device is invalid
	log.Printf("WARN: Publish event to %s but no publisher defined", key)
	return nil
}
