package api

import (
	"github.com/peter-mount/go-kernel/rest"
)

func (s *Service) getDevices(r *rest.Rest) error {
	r.Status(200).JSON().Value(s.model.GetModel().Device)
	return nil
}

func (s *Service) getDevice(r *rest.Rest) error {
	deviceName := r.Var("id")

	device := s.model.GetModel().GetDevice(deviceName)
	if device == nil {
		r.Status(404)
	} else {
		r.Status(200).JSON().Value(device)
	}

	return nil
}

func (s *Service) activateDevice(r *rest.Rest) error {
	deviceName := r.Var("id")
	action := r.Var("action")

	ok, err := s.activateDeviceImpl(deviceName, action)
	if err != nil {
		return err
	}
	if ok {
		// 202 Accepted as we have passed the activation request to RabbitMQ->Zigbee->Device(s)
		r.Status(202).JSON()
		return nil
	}

	r.Status(404)
	return nil
}

func (s *Service) activateDeviceImpl(deviceName, action string) (bool, error) {
	device := s.model.GetModel().GetDevice(deviceName)
	if device != nil {
		actions, ok := device.Action[action]
		if ok {
			for _, action := range actions {
				err := s.Send(deviceName+".set", action)
				if err != nil {
					return false, err
				}
			}
			return true, nil
		}
	}
	return false, nil
}
