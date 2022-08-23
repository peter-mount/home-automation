package state

import (
	"context"
	"encoding/json"
)

type State map[string]interface{}

func (s *State) GetString(k string) string {
	if s != nil {
		if v, exists := (*s)[k]; exists {
			if st, ok := v.(string); ok {
				return st
			}
			if st, ok := v.(bool); ok {
				if st {
					return "true"
				}
				return "false"
			}
		}
	}
	return ""
}

const (
	Key              = "state"
	PreviousStateKey = "previousState"
	ServiceKey       = "state.service"
)

func GetState(ctx context.Context) *State {
	return ctx.Value(Key).(*State)
}

func GetPreviousState(ctx context.Context) *State {
	return ctx.Value(PreviousStateKey).(*State)
}

func UnmarshalState(b []byte) (*State, error) {
	state := &State{}
	err := json.Unmarshal(b, state)
	if err != nil {
		return nil, err
	}
	return state, nil
}
