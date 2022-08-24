package enviro

import (
	"context"
	"encoding/json"
	"github.com/peter-mount/home-automation/util/graphite"
	mq2 "github.com/peter-mount/home-automation/util/mq"
	"log"
	"time"
)

type Enviro struct {
	mq        *mq2.MQ           `kernel:"inject"`
	queueName *mq2.Queue        `kernel:"config,enviroQueue"`
	graphite  graphite.Graphite `kernel:"inject"`
}

func (m *Enviro) Start() error {

	err := m.mq.ConsumeTask(m.queueName, "graphite", mq2.Guard(m.consume))
	if err != nil {
		return err
	}

	return err
}

func (m *Enviro) consume(ctx context.Context) error {
	body := mq2.Delivery(ctx)

	data := make(map[string]interface{})
	err := json.Unmarshal(body.Body, &data)
	if err != nil {
		log.Println(err)
		return err
	}

	// Timestamp in UTC of message
	ts, ok := data["timestamp"].(string)
	if !ok {
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05Z", ts+"Z")
	if err != nil {
		return err
	}

	routingKey := body.RoutingKey
	for k, v := range data {
		if !(k == "device" || k == "timestamp") {

			// TODO convert bool to 0 or 1 - place this in Publish if this works
			if b, ok := v.(bool); ok {
				if b {
					v = 1
				} else {
					v = 0
				}
			}

			err = m.graphite.Publish(t, routingKey+"."+k, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
