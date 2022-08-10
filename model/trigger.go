package model

import (
	"context"
	"github.com/peter-mount/home-automation/mq"
	"github.com/peter-mount/home-automation/state"
)

type Triggers []*Trigger

type Trigger struct {
	Device   string `yaml:"device"`
	Platform string `yaml:"platform"`
	From     string `yaml:"from,omitempty"`
	To       string `yaml:"to,omitempty"`
}

func (t Triggers) IsTriggered(ctx context.Context) bool {
	for _, trigger := range t {
		if trigger.IsFired(ctx) {
			return true
		}
	}
	return false
}

// IsFired returns true if the context contains data suitable to trigger this trigger
func (t *Trigger) IsFired(ctx context.Context) bool {
	msg := mq.Delivery(ctx)
	if state.FixKey(msg.RoutingKey) != state.FixKey(t.Device) {
		return false
	}

	switch t.Platform {
	case "action":
		return t.isAction(ctx)
	case "state":
		return t.isState(ctx)
	case "occupancy":
		return t.test(ctx, "occupancy")
	default:
		return false
	}
}

func (t *Trigger) isAction(ctx context.Context) bool {

	r := true
	if t.From != "" {
		r = r && state.GetPreviousState(ctx).GetString("action") == t.From
	}
	if t.To != "" {
		r = r && state.GetState(ctx).GetString("action") == t.To
	}
	return r
}

func (t *Trigger) isState(ctx context.Context) bool {
	r := true
	if t.From != "" {
		r = r && state.GetPreviousState(ctx).GetString("state") == t.From
	}
	if t.To != "" {
		r = r && state.GetState(ctx).GetString("state") == t.To
	}
	return r
}

func (t *Trigger) test(ctx context.Context, field string) bool {
	r := true
	if t.From != "" {
		r = r && state.GetPreviousState(ctx).GetString(field) == t.From
	}
	if t.To != "" {
		r = r && state.GetState(ctx).GetString(field) == t.To
	}
	return r
}
