package enviro

import (
	"context"
	"encoding/json"
	"github.com/peter-mount/home-automation/util/graphite"
	mq2 "github.com/peter-mount/home-automation/util/mq"
	"log"
	"strings"
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
	// The old format was "YYYY-MM-DD HH:MM:SS" whilst the new format is correct
	// So convert the old to new format for timestamp
	ts = strings.ReplaceAll(ts, " ", "T")
	if !strings.HasSuffix(ts, "Z") {
		ts = ts + "Z"
	}

	t, err := time.Parse("2006-01-02T15:04:05Z", ts)
	if err != nil {
		return err
	}

	return m.submitReadings(t, body.RoutingKey, data)
}

func (m *Enviro) submitReadings(t time.Time, routingKey string, data map[string]interface{}) error {
	for k, v := range data {
		switch k {
		// Ignore old format keys that are not metrics
		case "device":
		case "timestamp":
			// New format keys to ignore
		case "nickname":
		case "model":
		case "uid":
		// New format, readings are in this key
		case "readings":
			// Recurse if readings is a map
			if readings, ok := v.(map[string]interface{}); ok {
				err := m.submitReadings(t, routingKey, readings)
				if err != nil {
					return err
				}
			}
			// Old format, anything else is a metric
		default:
			err := m.submitReading(t, routingKey, k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Enviro) submitReading(t time.Time, routingKey, k string, v interface{}) error {
	// Convert a bool type to a 1 or 0 numeric type
	if b, ok := v.(bool); ok {
		if b {
			v = 1
		} else {
			v = 0
		}
	}

	return m.graphite.Publish(t, routingKey+"."+k, v)
}
