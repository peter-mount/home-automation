package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"strings"
	"time"
)

type Publisher struct {
	Exchange  string            `yaml:"exchange"`    // Exchange to submit to
	Key       string            `yaml:"DeliveryKey"` // Key to use for SimplePublish
	Mandatory bool              `yaml:"mandatory"`   // Publish mode
	Immediate bool              `yaml:"immediate"`   // Publish mode
	Replace   map[string]string `yaml:"replace"`     // Replace prefix table
	Ignore    []string          `yaml:"ignore"`      // Ignore prefixes
	Debug     bool              `yaml:"debug"`       // Debug mode
	Disabled  bool              `yaml:"disabled"`    // Publish disabled
	channel   *amqp.Channel
	mq        *MQ
}

func (p *Publisher) SimplePublish(msg []byte) error {
	return p.Post(p.Key, msg, nil, time.Now())
}

func (p *Publisher) Publish(key string, msg []byte) error {
	return p.Post(key, msg, nil, time.Now())
}

func (p *Publisher) PublishJSON(key string, payload interface{}) error {
	msg, err := json.Marshal(payload)
	if err == nil {
		err = p.Publish(key, msg)
	}
	return err
}

func (p *Publisher) Post(key string, body []byte, headers amqp.Table, timestamp time.Time) error {

	key = p.EncodeKey(key)

	// Check for ignored entries
	for _, v := range p.Ignore {
		if strings.HasPrefix(key, v+".") {
			if p.Debug {
				log.Printf("Ignoring %q:%q %s", p.Exchange, key, body)
			}
			return nil
		}
	}

	if p.Debug {
		log.Printf("Post %q:%q %s", p.Exchange, key, body)
	}

	if p.Disabled {
		return nil
	}

	return p.channel.Publish(
		p.Exchange,
		key,
		p.Mandatory,
		p.Immediate,
		amqp.Publishing{
			Headers:   headers,
			Timestamp: timestamp,
			Body:      body,
		},
	)
}

// EncodeKey converts any spaces or / in the DeliveryKey to . so they are more compatible with
// graphite. The output will also in lower case.
func (p *Publisher) EncodeKey(key string) string {
	key = EncodeKey(key)
	for k, v := range p.Replace {
		if strings.HasPrefix(key, k+".") {
			a := v + key[len(k):]
			key = a
		}
	}
	return key
}

func EncodeKey(key string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(key), " ", "."), "/", ".")
}
