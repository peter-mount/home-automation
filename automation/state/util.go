package state

import (
	"context"
	"strings"
)

// FixKey ensures that states use '.' as separators. '/' will be converted to '.'
func FixKey(k string) string {
	return strings.ReplaceAll(strings.ReplaceAll(k, "/", "."), " ", ".")
}

// GetService returns the Service instance within the context
func GetService(ctx context.Context) Service {
	return ctx.Value(ServiceKey).(Service)
}
