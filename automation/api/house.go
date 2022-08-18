package api

import "github.com/peter-mount/go-kernel/rest"

func (s *Service) getHouse(r *rest.Rest) error {
	r.Status(200).JSON().Value(s.house.GetModel())
	return nil
}
