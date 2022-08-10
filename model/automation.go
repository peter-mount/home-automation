package model

import (
	"context"
)

type Automation struct {
	Alias       string     `yaml:"alias,omitempty"`
	Description string     `yaml:"description"`
	Triggers    Triggers   `yaml:"triggers,omitempty"`
	Conditions  Conditions `yaml:"conditions,omitempty"`
	Actions     []*Action  `yaml:"actions,omitempty"`
}

func (a *Automation) IsTriggered(ctx context.Context) bool {
	return a.Triggers.IsTriggered(ctx) && a.Conditions.Match(ctx)
}
