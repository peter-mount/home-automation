package api

import (
	"encoding/json"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/home-automation"
	"github.com/peter-mount/home-automation/cache"
	"github.com/peter-mount/home-automation/mq"
	"time"
)

type Service struct {
	mq        *mq.MQ              `kernel:"inject"`
	publisher *mq.Publisher       `kernel:"config,apiPublisher"`
	model     *automation.Service `kernel:"inject"`
	rest      *rest.Server        `kernel:"inject"`
	cache     *cache.Cache        `kernel:"inject"`
}

func (s *Service) Start() error {

	err := s.mq.AttachPublisher(s.publisher)
	if err != nil {
		return err
	}

	s.rest.Handle("/api/zigbee/devices", s.listCacheDevices).Methods("GET")

	s.rest.Handle("/api/zigbee/{device:[0-9a-z/]+}/get", s.getCacheState).Methods("GET")

	s.rest.Handle("/api/zigbee/{device:[0-9a-z/]+}/set/{val}", s.setCacheState).Methods("GET")

	s.rest.Handle("/api/zigbee/{device:[0-9a-z/]+}/set", s.setCacheStateGeneric).Methods("POST")

	s.rest.Handle("/api/house", s.getHouse).Methods("GET")
	s.rest.Handle("/api/device", s.getDevices).Methods("GET")
	s.rest.Handle("/api/device/{id:[0-9a-z/.]+}", s.getDevice).Methods("GET")
	s.rest.Handle("/api/device/{id:[0-9a-z/.]+}:{action}", s.activateDevice).Methods("POST")
	s.rest.Handle("/api/scene", s.getScenes).Methods("GET")
	s.rest.Handle("/api/scene/{scene:[0-9a-z/]+}", s.getScene).Methods("GET")
	s.rest.Handle("/api/scene/{scene:[0-9a-z/]+}", s.activateScene).Methods("POST")

	return nil
}

// Send a message to zigbee2mqtt.
// []byte and string are sent as-is otherwise the message is marshaled into JSON before sending.
func (s *Service) Send(device string, msg interface{}) error {
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

	return s.publisher.Post(device, data, nil, time.Now())
}

func (s *Service) getHouse(r *rest.Rest) error {
	r.Status(200).JSON().Value(s.model.GetModel())
	return nil
}
