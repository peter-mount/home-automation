package model

import "github.com/peter-mount/home-automation/state"

type House struct {
	Device     map[string]*Device     `yaml:"device"`
	Scene      map[string]*Scene      `yaml:"scene"`
	Automation map[string]*Automation `yaml:"automation"`
}

func (h *House) GetDevice(k string) *Device {
	return h.Device[k]
}

func (h *House) GetScene(k string) *Scene {
	return h.Scene[state.FixKey(k)]
}
