package model

import (
	"context"
	"github.com/peter-mount/home-automation/automation/state"
	"strings"
)

type Conditions struct {
	Type       string       `yaml:"type"`
	Conditions []*Condition `yaml:"condition"`
}

type Condition struct {
	Type   string `yaml:"type"`             // Type of condition
	Device string `yaml:"device,omitempty"` // Device to inspect
	State  string `yaml:"state,omitempty"`  // State to inspect for Device
	Global string `yaml:"global,omitempty"` // Global variable to inspect
	Value  string `yaml:"value"`            // Value to test against for Global
}

func (c Conditions) Match(ctx context.Context) bool {
	switch strings.ToLower(c.Type) {
	case "or":
		// Pass fast on first match
		for _, cond := range c.Conditions {
			if cond.Matches(ctx) {
				return true
			}
		}
		return false

		// Default is "and"
	default:
		// Fail fast on first non-match
		for _, cond := range c.Conditions {
			if !cond.Matches(ctx) {
				return false
			}
		}
		return true
	}
}

func (c *Condition) Matches(ctx context.Context) bool {

	// Get values to compare
	var a, b string
	switch {
	// Device state
	case c.State != "":
		a = state.GetService(ctx).GetState(c.Device).GetString("state")
		b = c.State

	case c.Global != "":
		a = state.GetService(ctx).GetGlobal(c.Global)
		b = c.Value

	default:
		return false
	}

	a = strings.ToUpper(a)
	b = strings.ToUpper(b)

	r := false
	switch c.Type {
	case "equals":
		r = a == b

	case "notEquals":
		r = a != b

	}
	return r
}
