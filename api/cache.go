package api

import (
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/home-automation/util"
	"io/ioutil"
	"strings"
)

func (s *Service) listCacheDevices(r *rest.Rest) error {
	return util.DoRequest(r, func(status *util.Status) error {
		var r []*util.Event

		for _, d := range s.cache.GetDevices() {
			n := d.FriendlyName
			if !strings.HasPrefix(n, "zigbee2mqtt/") {
				n = "zigbee2mqtt/" + n
			}

			s, exists := s.cache.GetState(n)
			if !exists {
				s = nil
			}
			r = append(r, util.NewEvent(n, "status", s, d))
		}

		status.Data = r

		return nil
	})
}

// Get current device state (as in local cache)
func (s *Service) getCacheState(r *rest.Rest) error {
	return util.DoRequest(r, func(status *util.Status) error {
		device := r.Var("device")
		state, exists := s.cache.GetState(device)

		if exists {
			status.Data = util.NewEvent(device, "status", state, s.cache.GetDevice(device))
		} else {
			status.Status = 404
		}

		return nil
	})
}

func (s *Service) setCacheState(r *rest.Rest) error {
	return util.DoRequest(r, func(status *util.Status) error {
		device := r.Var("device")

		// Even if it doesn't exist, allow it as it might be a new device
		// or one that's not reported since the service started
		state, _ := s.cache.GetState(device)
		stateUpdate := state.DoUpdate()
		stateUpdate.SetState(r.Var("val"))
		status.Data = s

		return s.publisher.PublishApi(device+".set", stateUpdate)
	})
}

func (s *Service) setCacheStateGeneric(r *rest.Rest) error {
	return util.DoRequest(r, func(status *util.Status) error {
		device := r.Var("device")

		rdr, err := r.BodyReader()
		if err != nil {
			return err
		}

		payload, err := ioutil.ReadAll(rdr)
		if err != nil {
			return err
		}

		return s.publisher.PublishApi(device+".set", payload)
	})
}
