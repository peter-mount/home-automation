package util

import "github.com/peter-mount/go-kernel/rest"

type Status struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func DoRequest(r *rest.Rest, f func(*Status) error) error {
	var s Status
	err := f(&s)
	if err != nil {
		// Force 500 as this is a server error
		s.Status = 500
		s.Message = err.Error()
		s.Data = err
	} else if s.Status == 0 {
		// If s.Status is not 0 then leave it as the handler has set it
		s.Status = 200
	}
	r.Status(s.Status).JSON().Value(s)
	return nil
}
