package mq

import (
	"github.com/streadway/amqp"
	"log"
)

type Queue struct {
	Name       string    `yaml:"name"`
	Binding    []Binding `yaml:"binding"`
	Durable    bool      `yaml:"durable"`
	AutoDelete bool      `yaml:"autoDelete"`
	channel    *amqp.Channel
	mq         *MQ
}

type Binding struct {
	Topic string `yaml:"topic"`
	Key   string `yaml:"key"`
}

func (q *Queue) process(ch <-chan amqp.Delivery, f Task) {
	for {
		msg := <-ch
		q.logError(f(msg))
	}
}

func (q *Queue) logError(err error) {
	if err != nil {
		log.Printf("error, queue=%q, error=%v", q.Name, err)
	}
}
