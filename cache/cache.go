package cache

import (
	"context"
	"encoding/json"
	"github.com/peter-mount/go-kernel/util/task"
	"github.com/peter-mount/home-automation/mq"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

// Cache implements a service which stores the available devices and their current states
type Cache struct {
	mq        *mq.MQ        `kernel:"inject"`
	queueName *mq.Queue     `kernel:"config,bridgeQueue"`
	publisher *mq.Publisher `kernel:"config,bridgePublisher"`
	worker    task.Queue    `kernel:"worker"`
	mutex     sync.Mutex
	devices   map[string]*Device // Map of devices
	state     map[string]*State  // Map of current state
}

func (c *Cache) Start() error {
	c.devices = map[string]*Device{}
	c.state = map[string]*State{}

	err := c.mq.AttachPublisher(c.publisher)
	if err != nil {
		return err
	}

	err = c.mq.ConsumeTask(c.queueName, "graphite", mq.Guard(c.updateCache))
	if err != nil {
		return err
	}

	c.worker.AddTask(c.refresh)

	return nil
}

// refresh requests data from zigbee2mqtt.
// It's done as a worker task as it's only requested once the system is up and running
func (c *Cache) refresh(_ context.Context) error {
	log.Println("Requesting state from zigbee2mqtt")

	_ = c.Send("zigbee2mqtt/bridge/config/devices/get", "")

	return nil
}

func (c *Cache) GetState(device string) (*State, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	device = c.publisher.EncodeKey(device)
	s, exists := c.state[device]
	if !exists {
		s = &State{}
	}
	return s, exists
}

func (c *Cache) SetState(device string, state *State) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	device = c.publisher.EncodeKey(device)
	c.state[device] = state
}

func (c *Cache) updateCache(ctx context.Context) error {

	msg := mq.Delivery(ctx)
	switch msg.RoutingKey {
	case "zigbee2mqtt.bridge.log":
		c.log(msg)
	case "zigbee2mqtt.bridge.logging":
		c.log(msg)

	case "zigbee2mqtt.bridge.devices":
		//fmt.Printf("%s\n", msg.Body)

	default:
		if !(strings.HasSuffix(msg.RoutingKey, ".get") || strings.HasSuffix(msg.RoutingKey, ".set")) {
			s, _ := c.GetState(msg.RoutingKey)

			err := json.Unmarshal(msg.Body, s)
			if err == nil {
				c.SetState(msg.RoutingKey, s)
				/*if !exists {
				  log.Printf("New device %q %s", msg.RoutingKey, string(msg.Body))
				}*/
			}
		}
	}

	return nil
}

// Send a message to zigbee2mqtt.
// []byte and string are sent as-is otherwise the message is marshaled into JSON before sending.
func (c *Cache) Send(device string, msg interface{}) error {
	var data []byte

	if b, ok := msg.([]byte); ok {
		data = b
	} else if s, ok := msg.(string); ok {
		data = []byte(s)
	} else {
		b, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		data = b
	}

	return c.publisher.Post(device, data, nil, time.Now())
}

func (c *Cache) GetDevice(name string) *Device {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//name = b.publisher.EncodeKey(name)
	if strings.HasPrefix(name, "zigbee2mqtt/") {
		name = name[12:]
	}
	log.Println(name)
	if d, exists := c.devices[name]; exists {
		return d
	}
	return nil
}

func (c *Cache) addDevice(d *Device) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if d.FriendlyName != "" {
		c.devices[d.FriendlyName] = d
	}
}

// GetDevices returns a list of devices
func (c *Cache) GetDevices() []*Device {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var r []*Device
	for _, v := range c.devices {
		r = append(r, v)
	}

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].FriendlyName < r[j].FriendlyName
	})

	return r
}
