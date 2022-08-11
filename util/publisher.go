package util

import (
	"github.com/peter-mount/home-automation/mq"
)

// Publisher handles the publishing of messages back to RabbitMQ using the
// config defined under the "automationPublisher" key.
// Most services will use this publisher for everything sent to RabbitMQ.
type Publisher struct {
	mq        *mq.MQ        `kernel:"inject"`
	publisher *mq.Publisher `kernel:"config,automationPublisher"`
}

func (s *Publisher) Start() error {
	err := s.mq.AttachPublisher(s.publisher)
	if err != nil {
		return err
	}

	return nil
}

// PublishJSON sends the payload as a JSON object using the supplied routing key
func (s *Publisher) PublishJSON(key string, payload interface{}) error {
	return s.publisher.PublishJSON(key, payload)
}

// PublishApi sends the payload using the supplied routing key.
// []byte and string are sent as-is otherwise the message is marshaled into JSON before sending.
func (s *Publisher) PublishApi(device string, msg interface{}) error {
	return s.publisher.PublishApi(device, msg)
}
