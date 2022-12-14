package cache

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Log struct {
	Level   string    `json:"level,omitempty"`
	Type    string    `json:"type"`
	Message []*Device `json:"message"`
}

func (c *cache) log(msg amqp.Delivery) {
	var l Log
	err := json.Unmarshal(msg.Body, &l)
	if err == nil {
		if l.Type == "devices" {
			c.logDevice(l)
			//fmt.Printf("%s\n", msg.Body)
		}
	}
}

func (c *cache) logDevice(l Log) {
	for _, msg := range l.Message {
		device := msg
		c.worker.AddTask(func(_ context.Context) error {
			c.addDevice(device)
			return c.publisher.PublishApi("zigbee2mqtt."+device.FriendlyName+".get", "{}")
		})
	}
}
