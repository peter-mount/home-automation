package api

import "github.com/peter-mount/go-kernel/rest"

func (s *Service) getScenes(r *rest.Rest) error {
	r.Status(200).JSON().Value(s.house.GetModel().Scene)
	return nil
}

func (s *Service) getScene(r *rest.Rest) error {
	sceneName := r.Var("scene")

	scene := s.house.GetModel().GetScene(sceneName)
	if scene == nil {
		r.Status(404)
	} else {
		r.Status(200).JSON().Value(scene)
	}

	return nil
}

func (s *Service) activateScene(r *rest.Rest) error {
	sceneName := r.Var("scene")

	scene := s.house.GetModel().GetScene(sceneName)
	if scene != nil {
		for deviceName, action := range scene.Devices {
			if _, err := s.activateDeviceImpl(deviceName, action); err != nil {
				return err
			}
		}

		// 202 Accepted as we have passed the activation request to RabbitMQ->Zigbee->Device(s)
		r.Status(202).JSON()
	} else {
		r.Status(404)
	}
	return nil
}
