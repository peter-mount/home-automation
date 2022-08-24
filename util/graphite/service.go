package graphite

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/peter-mount/go-kernel"
	mq2 "github.com/peter-mount/home-automation/util/mq"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Graphite handles receiving events from rabbit and logging the responses to graphite
// via RabbitMQ.
type Graphite interface {
	Publish(t time.Time, k string, v interface{}) error
}

func init() {
	kernel.RegisterAPI((*Graphite)(nil), &graphite{})
}

type graphite struct {
	mq        *mq2.MQ        `kernel:"inject"`
	queueName *mq2.Queue     `kernel:"config,graphiteQueue"`
	publisher *mq2.Publisher `kernel:"config,graphitePublisher"`
}

func (m *graphite) Start() error {
	err := m.mq.AttachPublisher(m.publisher)
	/*
	   if err == nil {
	     err = m.mq.ConsumeTask(m.queueName, "graphite", mq.Guard(m.logMessage))
	*/
	return err
}

// logMessage receives a message from rabbitmq
func (m *graphite) logMessage(ctx context.Context) error {
	body := mq2.Delivery(ctx)

	// Ignore bridge specific messages as this service deals with devices
	if strings.HasPrefix(body.RoutingKey, "zigbee2mqtt.bridge") {
		return nil
	}

	data := make(map[string]interface{})
	err := json.Unmarshal(body.Body, &data)
	if err != nil {
		log.Println(err)
		return err
	}

	// If last_seen is missing then ignore this message.
	// It'll probably be a message to zigbee2mqtt not from it.
	ts, ok := data["last_seen"].(string)
	if !ok {
		return nil
	}

	t, err := time.Parse("2006-01-02T15:04:05.999Z", ts)
	if err != nil {
		return err
	}

	for k, v := range data {
		err = m.Publish(t, body.RoutingKey+"."+k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *graphite) Publish(t time.Time, k string, v interface{}) error {
	var val string

	if v != nil {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Int:
			val = strconv.Itoa(v.(int))
		case reflect.Float64:
			val = fmt.Sprintf("%.3f", v.(float64))
		case reflect.String:
			switch strings.ToLower(v.(string)) {
			case "on":
				val = "1"
			case "off":
				val = "0"
			case "true":
				val = "1"
			case "false":
				val = "0"
			}
		case reflect.Bool:
			if v.(bool) {
				val = "1"
			} else {
				val = "0"
			}
		}
	}

	if val != "" {
		ts := t.UTC()
		key := m.publisher.EncodeKey(k)
		msg := fmt.Sprintf("%s %s %d", key, val, ts.Unix())

		return m.publisher.Post(key, []byte(msg), nil, ts)
	}

	return nil
}
