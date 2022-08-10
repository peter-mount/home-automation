package model

import "github.com/peter-mount/home-automation/state"

type Device struct {
	Description string                    `yaml:"description"`
	Action      map[string][]*state.State `yaml:"action"`
}
