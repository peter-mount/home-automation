package api

import (
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/home-automation/cache"
	"github.com/peter-mount/home-automation/model"
	"github.com/peter-mount/home-automation/mq"
	automation "github.com/peter-mount/home-automation/util"
)

type Service struct {
	mq        *mq.MQ                `kernel:"inject"`
	publisher *automation.Publisher `kernel:"inject"`
	house     *model.Service        `kernel:"inject"`
	rest      *rest.Server          `kernel:"inject"`
	cache     *cache.Cache          `kernel:"inject"`
}

func (s *Service) Start() error {

	s.rest.Handle("/api/cache/devices", s.listCacheDevices).Methods("GET")
	s.rest.Handle("/api/cache/{device:[0-9a-z/]+}/get", s.getCacheState).Methods("GET")
	s.rest.Handle("/api/cache/{device:[0-9a-z/]+}/set/{val}", s.setCacheState).Methods("GET")
	s.rest.Handle("/api/cache/{device:[0-9a-z/]+}/set", s.setCacheStateGeneric).Methods("POST")

	s.rest.Handle("/api/house", s.getHouse).Methods("GET")

	s.rest.Handle("/api/device", s.getDevices).Methods("GET")
	s.rest.Handle("/api/device/{id:[0-9a-z/.]+}", s.getDevice).Methods("GET")
	s.rest.Handle("/api/device/{id:[0-9a-z/.]+}:{action}", s.activateDevice).Methods("POST")

	s.rest.Handle("/api/scene", s.getScenes).Methods("GET")
	s.rest.Handle("/api/scene/{scene:[0-9a-z/]+}", s.getScene).Methods("GET")
	s.rest.Handle("/api/scene/{scene:[0-9a-z/]+}", s.activateScene).Methods("POST")

	return nil
}
